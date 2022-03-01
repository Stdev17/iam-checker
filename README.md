# IAM Lifecycle Checker (인터뷰용 과제)

해당 AWS 계정의 IAM User들의 Access Key Pair 생성 시기를 조회하여, 최근 N시간 이내에 이를 교체하지 않은 IAM User를 찾아 출력 합니다.

### 구현 사양

- 해당 IAM User가 가진 **모든** 만료된 Access Key Pair를 출력합니다.
  - 결과 파일은 .json 형태로 /tmp 디렉토리에 저장됩니다. 디버깅을 위해 에러가 없을 때도 로그를 출력하도록 설정해 두었습니다. 컨테이너 환경에서 저 파일은 ephemeral하므로, 실무에서는 PersistentVolume에 마운트하거나, MQ 등으로 결과 값을 중앙 로그 파이프라인으로 전송하는 practice가 일반적입니다.
- 환경변수의 LIFETIME 값을 기준으로 동작하며, 환경변수는 쿠버네티스 Secret 파일 형태로 Pod/CronJob 생성 시에 주입됩니다.
- 주기적으로 Access Key Pair 생성 시기를 검사한다는 시나리오 하에, CronJob으로 매주 월요일 오전 10시에 작업을 수행한다는 설정을 부여했습니다. Pod은 일회성으로 실행되고 terminated된 이후 재실행되지 않도록 합니다. (restartPolicy: \"Never\")

### 고민해 볼 것들

- 클라우드 VM의 실행 환경은 amd64 linux가 de facto standard인 현황입니다. 그런데 서버에서 동작하는 애플리케이션 외에 사용자가 직접 쓰는 사내 툴이라면 환경이 꽤 다양해질 것입니다. (amd64 vs. arm64, windows/darwin/linux) 멀티 플랫폼 환경에서 별도의 포팅 작업이 필요치 않도록 Go에서 제공하는 versatile한 패키지(os, path/filepath)를 활용하여 코드를 작성해 보고 있습니다만, 실제로 환경별 바이너리를 빌드, 배포할 때 품을 줄일 만한 루틴을 더 알아나가야 하겠습니다.
  - Makefile 등을 활용한 개발-빌드-배포 이터레이션에 아직 익숙치 못한데, 프로덕션 CICD 파이프라인에서 각 스택이 어떻게 접목되는지에 대한 이해도가 장차 선행되어야 할 것입니다.
- Vault 같은 Secret Store를 자체 Secret 대신 활용하여 개발 편의성과 보안이 향상될지 리서치해 볼 가치가 있겠습니다.
- 환경변수 대신 url query string을 받아 동작한다면 HTTP GET 요청을 listen하는 웹서버 구현이 필요할 것입니다. 이때는 하나의 서비스로 간주하여, Deployment 단위로 추상화하여 다루는 것이 일반적이겠습니다.
  - 주기적으로 로그 자체를 기록한다면 CronJob 식 접근이, 다른 서비스에서 결과값을 제공받아 후속 처리(SNS로 alert나 Access Key Pair 삭제 등)를 진행한다면 서비스 형태로 띄워 로직을 디커플링하는 접근이 유리하겠습니다.
