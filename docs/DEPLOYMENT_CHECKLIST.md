# Deployment Checklist

## Pre-Deployment

### Code
- [ ] All tests passing
- [ ] Code reviewed and approved
- [ ] No hardcoded secrets
- [ ] Environment variables documented
- [ ] Database migrations tested

### Security
- [ ] JWT secrets rotated (production)
- [ ] SSL certificates valid
- [ ] CORS configured for production domains
- [ ] Rate limiting enabled
- [ ] RBAC permissions verified

### Infrastructure
- [ ] Database backups configured
- [ ] Redis persistence enabled
- [ ] Object storage (MinIO) with retention policy
- [ ] Monitoring and alerting set up
- [ ] Log aggregation configured

## Deployment Steps

### 1. Database Migration
```bash
# Run migrations before deploying new code
docker compose exec api migrate up
```

### 2. Pull Latest Images
```bash
docker compose pull
```

### 3. Deploy Backend
```bash
docker compose up -d api worker
# Wait for health checks
docker compose ps
```

### 4. Deploy Frontend
```bash
docker compose up -d frontend nginx
```

### 5. Verify Deployment
```bash
# Check API health
curl https://api.pmapp.com/health

# Check frontend
curl https://pmapp.com

# Check logs
docker compose logs --tail=100 api
```

## Post-Deployment

### Verification
- [ ] Login flow works
- [ ] Create project successful
- [ ] Create task successful
- [ ] Real-time updates working
- [ ] File upload working
- [ ] Email notifications sending

### Monitoring
- [ ] Check Prometheus metrics
- [ ] Verify no error spikes in logs
- [ ] Check database connection pool
- [ ] Verify Redis cache hit rate

### Rollback Plan
```bash
# Rollback to previous version
docker compose pull
docker compose up -d api worker

# If database migration needed rollback:
docker compose exec api migrate down
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| DATABASE_URL | PostgreSQL connection string | Yes |
| REDIS_URL | Redis connection string | Yes |
| JWT_SECRET | Access token secret (min 32 chars) | Yes |
| JWT_REFRESH_SECRET | Refresh token secret | Yes |
| MINIO_ENDPOINT | S3/MinIO endpoint | Yes |
| MINIO_ACCESS_KEY | Storage access key | Yes |
| MINIO_SECRET_KEY | Storage secret key | Yes |
| AI_API_KEY | OpenAI/Anthropic API key | No |
| GRAFANA_PASSWORD | Grafana admin password | No |

## Health Check Endpoints

| Service | Endpoint | Expected |
|---------|----------|----------|
| API | GET /health | 200 OK |
| Frontend | GET / | 200 OK |
| MinIO | GET /minio/health/live | 200 OK |
| PostgreSQL | pg_isready | 0 |
| Redis | redis-cli ping | PONG |
