import { useEffect, useRef, useState } from 'react';
import type { Camera } from '../features/api/apiSlice';

// Определим интерфейс для пропсов хука
interface UseWebRTCStreamProps {
  camera: Camera | null;
}

export const useWebRTCStream = ({ camera }: UseWebRTCStreamProps) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const peerConnectionRef = useRef<RTCPeerConnection | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!camera) return;

    const streamUrl = `ws://localhost:8889/local_cam`;

    const pc = new RTCPeerConnection({
      iceServers: [{ urls: 'stun:stun.l.google.com:19302' }],
    });
    peerConnectionRef.current = pc;

    pc.onicecandidate = (event) => {
      if (event.candidate) {
        // Отправляем ICE candidate на сервер
        ws.send(JSON.stringify({ type: 'ice', candidate: event.candidate }));
      }
    };

    pc.ontrack = (event) => {
      if (videoRef.current) {
        videoRef.current.srcObject = event.streams[0];
        setIsConnected(true);
        setError(null);
      }
    };

    const ws = new WebSocket(streamUrl);

    ws.onopen = async () => {
      console.log('WebSocket connection opened');
      pc.addTransceiver('video', { direction: 'recvonly' });
      const offer = await pc.createOffer();
      await pc.setLocalDescription(offer);
      ws.send(JSON.stringify({ type: 'offer', sdp: pc.localDescription }));
    };

    ws.onmessage = async (event) => {
      try {
        const message = JSON.parse(event.data);

        if (message.type === 'answer') {
          const remoteDesc = new RTCSessionDescription(message.sdp);
          await pc.setRemoteDescription(remoteDesc);
        } else if (message.type === 'ice') {
          await pc.addIceCandidate(new RTCIceCandidate(message.candidate));
        }
      } catch (e) {
        console.error('Failed to parse message or set description:', e);
        setError('Failed to establish WebRTC connection.');
      }
    };

    ws.onerror = (err) => {
      console.error('WebSocket error:', err);
      setError('WebSocket connection failed.');
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed');
      setIsConnected(false);
    };

    return () => {
      ws.close();
      pc.close();
    };
  }, [camera]);

  return { videoRef, isConnected, error };
};
