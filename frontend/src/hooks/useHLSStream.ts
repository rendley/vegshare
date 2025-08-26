import { useEffect, useRef } from 'react';
import Hls from 'hls.js';
import type { Camera } from '../features/api/apiSlice';

interface UseHLSStreamProps {
  camera: Camera | null;
}

export const useHLSStream = ({ camera }: UseHLSStreamProps) => {
  const videoRef = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    if (!camera || !videoRef.current) return;

    // HLS URL, который предоставляет mediamtx
    const hlsUrl = `/api/v1/stream/hls/${camera.rtsp_path_name}/index.m3u8`;

    const videoElement = videoRef.current;
    let hls: Hls | null = null;

    if (Hls.isSupported()) {
      console.log("Using hls.js for playback");
      hls = new Hls();
      hls.loadSource(hlsUrl);
      hls.attachMedia(videoElement);
      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        videoElement.play().catch(e => console.error("Autoplay failed", e));
      });
    } else if (videoElement.canPlayType('application/vnd.apple.mpegurl')) {
      console.log("Using native HLS support");
      // Для Safari, который поддерживает HLS нативно
      videoElement.src = hlsUrl;
      videoElement.addEventListener('loadedmetadata', () => {
        videoElement.play().catch(e => console.error("Autoplay failed", e));
      });
    }

    return () => {
      if (hls) {
        hls.destroy();
      }
    };
  }, [camera]);

  return { videoRef, isConnected: true, error: null }; // Упрощаем, считаем что всегда подключено
};