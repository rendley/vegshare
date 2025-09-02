import React from 'react';
import { useHLSStream } from '../hooks/useHLSStream';
import type { Camera } from '../features/api/apiSlice';
import { Card, CardContent, CardHeader, Typography, Box } from '@mui/material';

interface VideoPlayerProps {
  camera: Camera;
}

const VideoPlayer: React.FC<VideoPlayerProps> = ({ camera }) => {
  const { videoRef, isConnected, error } = useHLSStream({ camera });

  return (
    <Card sx={{ maxWidth: 720, m: 1 }}>
      <CardHeader 
        title={camera.name} 
        subheader={`Status: ${isConnected ? 'Connected' : 'Disconnected'}`} 
      />
      <CardContent sx={{ pt: 0 }}>
        <Box sx={{ backgroundColor: '#000', borderRadius: 1 }}>
          <video ref={videoRef} style={{ width: '100%', display: 'block' }} autoPlay playsInline muted controls />
        </Box>
        {error && (
          <Typography variant="caption" color="error" sx={{ mt: 1, display: 'block' }}>
            {error}
          </Typography>
        )}
      </CardContent>
    </Card>
  );
};

export default VideoPlayer;
