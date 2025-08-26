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

    const token = localStorage.getItem('token');
    if (!token) {
      console.error("Auth token not found for HLS stream");
      return;
    }

    // URL для основного плейлиста. Токен добавляется сразу.
    const hlsUrl = `/api/v1/stream/hls/${camera.rtsp_path_name}/index.m3u8?token=${token}`;

    const videoElement = videoRef.current;
    let hls: Hls | null = null;

    // Конфигурация для hls.js
    const hlsConfig = {
      loader: class CustomLoader extends Hls.DefaultConfig.loader {
        constructor(config: any) {
          super(config);
          const oldLoad = this.load.bind(this);
          this.load = (context, config, callbacks) => {
            // Используем type guard "in" для проверки наличия свойства frag.
            // Это доказывает TypeScript, что мы работаем с контекстом сегмента.
            if ('frag' in context) {
              // Добавляем токен к URL сегмента.
              // Важно использовать `&`, так как `?` уже используется для токена в основном URL.
              context.url += `?token=${token}`;
            }
            oldLoad(context, config, callbacks);
          };
        }
      },
    };

    if (Hls.isSupported()) {
      console.log("Using hls.js for playback with custom loader");
      hls = new Hls(hlsConfig);
      hls.loadSource(hlsUrl);
      hls.attachMedia(videoElement);
      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        videoElement.play().catch(e => console.error("Autoplay failed", e));
      });
    } else if (videoElement.canPlayType('application/vnd.apple.mpegurl')) {
      console.log("Using native HLS support (token auth for segments may not work)");
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

  return { videoRef, isConnected: true, error: null };
};