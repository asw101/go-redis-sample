export REDIS_PASSWORD=$(kubectl get secret --namespace "default" my-redis-redis-cluster -o jsonpath="{.data.redis-password}" | base64 --decode)

export REDIS_ADDR=$(kubectl get endpoints my-redis-redis-cluster -o=jsonpath='{.subsets[0].addresses[*].ip}')

kubectl run --namespace default go-redis-sample-produce \
  --rm --tty -i --restart='Never' \
  --env REDIS_PASSWORD="$REDIS_PASSWORD" \
  --env REDIS_ADDR="$REDIS_ADDR" \
  --image ghcr.io/asw101/go-redis-sample -- /app produce

