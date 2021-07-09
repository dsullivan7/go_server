aws ecs update-service \
  --cluster tchr-voice-cluster \
  --service tchr-voice-service-app \
  --task-definition tchr-voice-app:21 \
  --force-new-deployment
