#!/bin/bash

# Simple deployment script for MIJKomp Service
# Usage: ./deploy.sh [server_user@server_ip] [deploy_path]
# Example: ./deploy.sh root@192.168.1.100 /opt/mijkomp

set -e

SERVER=${1:-"root@localhost"}
DEPLOY_PATH=${2:-"/opt/mijkomp"}
SERVICE_NAME="mijkomp"

echo "Deploying MIJKomp Service to $SERVER:$DEPLOY_PATH"

# Build the application first
echo "Building application..."
./build-linux.sh

# Create deployment directory on server
echo "Creating deployment directory..."
ssh $SERVER "mkdir -p $DEPLOY_PATH"
ssh $SERVER "mkdir -p $DEPLOY_PATH/logs"
ssh $SERVER "mkdir -p $DEPLOY_PATH/uploads"

# Stop existing service if running
echo "Stopping existing service..."
ssh $SERVER "pkill -f mijkomp-service || true"

# Copy files to server
echo "Copying files to server..."
scp mijkomp-service $SERVER:$DEPLOY_PATH/
scp .env.example $SERVER:$DEPLOY_PATH/
scp -r docs/ $SERVER:$DEPLOY_PATH/

# Set permissions
echo "Setting permissions..."
ssh $SERVER "chmod +x $DEPLOY_PATH/mijkomp-service"
ssh $SERVER "chmod -R 755 $DEPLOY_PATH/docs"
ssh $SERVER "chmod -R 755 $DEPLOY_PATH/logs"
ssh $SERVER "chmod -R 755 $DEPLOY_PATH/uploads"

# Create .env if not exists
ssh $SERVER "if [ ! -f $DEPLOY_PATH/.env ]; then cp $DEPLOY_PATH/.env.example $DEPLOY_PATH/.env; echo 'Created .env file from template. Please configure it.'; fi"

# Start service in background
echo "Starting service..."
ssh $SERVER "cd $DEPLOY_PATH && nohup ./mijkomp-service > logs/app.log 2>&1 & echo \$! > mijkomp.pid"

# Wait a moment and check if service started
sleep 3
echo "Checking service status..."
if ssh $SERVER "ps -p \$(cat $DEPLOY_PATH/mijkomp.pid 2>/dev/null) > /dev/null 2>&1"; then
    echo "✅ Service started successfully!"
    echo "📋 Service PID: $(ssh $SERVER "cat $DEPLOY_PATH/mijkomp.pid 2>/dev/null || echo 'Unknown'")"
    echo "📁 Deployment path: $DEPLOY_PATH"
    echo "📊 Check logs: ssh $SERVER 'tail -f $DEPLOY_PATH/logs/app.log'"
    echo "🔍 Health check: curl http://$SERVER:5000/health"
else
    echo "❌ Service failed to start. Check logs:"
    ssh $SERVER "tail -20 $DEPLOY_PATH/logs/app.log"
    exit 1
fi

echo ""
echo "Deployment completed!"
echo "Next steps:"
echo "1. Configure $DEPLOY_PATH/.env with your database settings"
echo "2. Restart service: ssh $SERVER 'cd $DEPLOY_PATH && pkill -f mijkomp-service && nohup ./mijkomp-service > logs/app.log 2>&1 & echo \$! > mijkomp.pid'"
echo "3. Test API: curl http://$SERVER:5000/health"