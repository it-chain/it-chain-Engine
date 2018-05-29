---
본 문서는 peer 개발 시 소통을 위해 임의로 작성되었으며
전체 개발의 진행방향과 무관하게 소통을 위한 문서이기에 혼동이 없으셨으면 합니다.
-frontalnh(namhoon)
---

# 최초 노드 boot 시의 시나리오
최초의 노드를 bootnode로 하는 최초의 p2p 네트워크를 형성하는 시나리오이다.
이 경우 it-chain의 bootnode가 반드시 해당 pc 의 ip로 설정되어 있어야 한다.

개략적인 프로세스는 다음과 같다.

1. 추가될 노드와의 통신을 위한 준비
2. 노드 내 컴포넌트간 통신을 위한 AMQP 준비
3. p2p 네트워크 초기 세팅

## 추가될 노드와의 통신을 위한 준비
현재 boot 된 최초의 노드를 bootnode로 하여 최초의 p2p 네트워크를 형성하기 위해 다른 노드와 gRPC를 사용한 통신망을 구축하며, 실제적인 gateway 컴포넌트에서 이루어 진다. gateway 패키지의 `start()` 함수호출을 통해 gRPC interface를 구축한다.

이러한 함수호출을 통해 다음과 같은 작업들이 수행되게 된다.
1. gRPC interface 구축
2. amqp server 환경 세팅

이 과정에서 gRPC 인터페이스를 구축함에 있어 grpc lib을 사용하는 것이 일반적이지만, it-chain 에서는 이를 보다 쉽게 구현하기 위해 별도로 p2p 네트워크 라이브러리인 `bifrost` 라이브러리를 자체 제작하여 사용하고있다. `bifrost` 를 통해 gRPC 인터페이스를 보다 쉽게 구축 할 수 있을 뿐만 아니라 해당 p2p 네트워크의 전체 연결정보를 쉽게 관리하고 정보를 열람할 수 있다.

## 노드 내 컴포넌트간 통신을 위한 AMQP 세팅
노드 내의 각 컴포넌트들은 AMQP 를 사용하여 통신하며 rabbitmq 를 통해 구현되며, 이를 위해 amqp 서버를 구동이 필요하다.
amqp 서버 구동의 경우 앞의 gRPC 인터페이스를 구축함에 있어 gateway 패키지의 start() 함수를 호출하는 시점에서 amqp 서버의 구동까지 같이 수행하므로 별도로 다른 조작은 필요하지 않다.

**어느 단계에서 gateway.Start() 가 수행이 될지 정해지지 않았다.**


## peer 를 생성하고 leader를 선언
cli 에서 dial 하는 시점에서 bootnode ip를 입력해 주며, 이 과정에서 해당 boot node를 리더로 선출한다?!

이러한 leader의 정보는 peer의 peer table의 형태로 leveldb에 저장된다.



# 두번째 노드가 네트워크에 연결되는 시나리오
자신이 부트노드가 아니므로
1. 피어 생성
2. 다른 피어에게 알림
3. 기존 피어목록 생신
4. 내 피어 생성

---
아래 내용은 소통을 위해 생각의 흐름대로 임의로 작성하였으며
추후 개선될 예정입니다.
---


# 트랜잭션 발생 시나리오
1. 사용자가 tx 생성 요청
2. txpool 에 저장
3. txpool에서 consume
4. pool에 등록
5. api 호출
createtx
6. amqp 등록
7. create event 객체생성
8. pool에서 tx 생성
9. leveldb 저장


# 블록 생성
# 합의과정
# 블록 저장
# 블록 전파