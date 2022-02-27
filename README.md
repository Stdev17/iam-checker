# IAM Lifecycle Checker

해당 AWS 계정의 IAM User들의 Access Key Pair 생성시기를 조회하여, 최근 N시간 이내에 이를 교체하지 않은 IAM User를 찾아 출력 합니다.

- 해당 IAM User가 가진 **모든** 만료된 Access Key Pair를 출력합니다.
- 환경변수는 local에서 .env 파일로 정의되어 빌드되며, 추후에 외부에서 환경변수를 주입한다든가 하는 practice로 확장하여 적용할 수 있겠습니다.