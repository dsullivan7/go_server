aws ecs run-task \
  --cluster tchr-voice-cluster \
  --task-definition tchr-voice-migrate:4 \
  --count 1 \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-05787f76d23b5ca06,subnet-01e1b735d2f7df441],securityGroups=[sg-0b88230d956daa441]}"
