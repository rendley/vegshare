import React from 'react';
import { useHLSStream } from '../hooks/useHLSStream';
import type { Camera } from '../features/api/apiSlice';

interface VideoPlayerProps {
  camera: Camera;
}

const VideoPlayer: React.FC<VideoPlayerProps> = ({ camera }) => {
  const { videoRef, isConnected, error } = useHLSStream({ camera });

  return (
    <div style={{ border: '1px solid #ccc', padding: '10px', margin: '10px' }}>
      <h4>{camera.name}</h4>
      <video ref={videoRef} style={{ width: '100%' }} autoPlay playsInline muted />
      <div>
        {error && <p style={{ color: 'red' }}>Error: {error}</p>}
        <p>Status: {isConnected ? 'Connected' : 'Disconnected'}</p>
      </div>
    </div>
  );
};

export default VideoPlayer;
