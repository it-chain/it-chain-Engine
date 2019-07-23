/*
 * Copyright 2018 DE-labtory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pbftfx

import (
	"github.com/DE-labtory/it-chain/common"
	"github.com/DE-labtory/it-chain/common/rabbitmq/pubsub"
	"github.com/DE-labtory/it-chain/conf"
	"github.com/DE-labtory/it-chain/consensus/pbft"
	"github.com/DE-labtory/it-chain/consensus/pbft/api"
	"github.com/DE-labtory/it-chain/consensus/pbft/infra/adapter"
	"github.com/DE-labtory/it-chain/consensus/pbft/infra/mem"
	"github.com/DE-labtory/iLogger"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewParliamentRepository,
		mem.NewStateRepository,
		NewElectionService,
		NewPropagateService,
		NewElectionApi,
		NewParliamentApi,
		NewStateApi,
		adapter.NewElectionCommandHandler,
		adapter.NewConnectionEventHandler,
		NewNewLeaderCommandHandler,
		adapter.NewLeaderEventHandler,
		NewStartConsensusCommandHandler,
		NewPbftMsgHandler,
	),
	fx.Invoke(
		RegisterPubsubHandlers,
	),
)

func NewElectionService(config *conf.Configuration, nodeId common.NodeID) *pbft.ElectionService {
	return pbft.NewElectionService(nodeId, 30, pbft.NORMAL, 0)
}

func NewParliamentRepository(config *conf.Configuration, nodeId common.NodeID) *mem.ParliamentRepository {
	parliament := pbft.NewParliament()
	parliament.AddRepresentative(pbft.NewRepresentative(nodeId))

	if config.Engine.BootstrapNodeAddress == "" {
		parliament.SetLeader(nodeId)
	}

	return mem.NewParliamentRepositoryWithParliament(parliament)
}

func NewPropagateService(service common.EventService) *pbft.PropagateService {
	return pbft.NewPropagateService(service)
}

func NewStateApi(config *conf.Configuration, propagateService *pbft.PropagateService, service common.EventService, paliamentrepository *mem.ParliamentRepository, stateRepository *mem.StateRepository, publisherId common.NodeID) *api.StateApi {
	return api.NewStateApi(publisherId, propagateService, service, paliamentrepository, stateRepository)
}

func NewElectionApi(electionService *pbft.ElectionService, parliamentRepository *mem.ParliamentRepository, eventService common.EventService) *api.ElectionApi {
	return api.NewElectionApi(electionService, parliamentRepository, eventService)
}

func NewParliamentApi(config *conf.Configuration, parliamentRepository *mem.ParliamentRepository, eventService common.EventService, nodeId common.NodeID) *api.ParliamentApi {
	return api.NewParliamentApi(nodeId, parliamentRepository, eventService)
}

func NewNewLeaderCommandHandler(parliamentApi *api.ParliamentApi) *adapter.LeaderCommandHandler {
	return adapter.NewLeaderCommandHandler(parliamentApi)
}

func NewStartConsensusCommandHandler(stateApi *api.StateApi) *adapter.StartConsensusCommandHandler {
	return adapter.NewStartConsensusCommandHandler(stateApi)
}

func NewPbftMsgHandler(stateApi *api.StateApi) *adapter.PbftMsgHandler {
	return adapter.NewPbftMsgHandler(stateApi)
}

func RegisterPubsubHandlers(subscriber *pubsub.TopicSubscriber, pbftMsgHandler *adapter.PbftMsgHandler, electionCommandHandler *adapter.ElectionCommandHandler, connectionEventHandler *adapter.ConnectionEventHandler, leaderCommandHandler *adapter.LeaderCommandHandler, leaderEventHandler *adapter.LeaderEventHandler, startConsensusHandler *adapter.StartConsensusCommandHandler) {
	iLogger.Infof(nil, "[Main] Consensus is starting")

	if err := subscriber.SubscribeTopic("message.receive", electionCommandHandler); err != nil {
		panic(err)
	}

	if err := subscriber.SubscribeTopic("connection.*", connectionEventHandler); err != nil {
		panic(err)
	}

	if err := subscriber.SubscribeTopic("message.receive", leaderCommandHandler); err != nil {
		panic(err)
	}

	if err := subscriber.SubscribeTopic("leader.deleted", leaderEventHandler); err != nil {
		panic(err)
	}

	if err := subscriber.SubscribeTopic("block.consent", startConsensusHandler); err != nil {
		panic(err)
	}

	if err := subscriber.SubscribeTopic("message.receive", pbftMsgHandler); err != nil {
		panic(err)
	}
}
