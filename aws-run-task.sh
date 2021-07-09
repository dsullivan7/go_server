aws ecs run-task \
  --cluster tchr-voice-cluster \
  --task-definition tchr-voice-migrate:6 \
  --count 1 \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-07c96383382d15162,subnet-0eea121729d78f36d],securityGroups=[sg-0e5b2cdefdfe35af9]}"
