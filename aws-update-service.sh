aws ecs update-service \
  --cluster tchr-voice-cluster \
  --service tchr-voice-service-app \
  --task-definition tchr-voice-app:19 \
  --force-new-deployment
