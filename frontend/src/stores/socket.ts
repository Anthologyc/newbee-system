import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useSocketStore = defineStore('socket', () => {
  const isConnected = ref(false);
  let ws: WebSocket | null = null;

  // 事件监听器
  type Callback = (payload: any) => void;
  const listeners: Record<string, Callback[]> = {};

  const connect = (userId: number) => {
    if (ws && ws.readyState === WebSocket.OPEN) return;

    // 确保地址正确
    const url = `ws://localhost:8080/ws?user_id=${userId}`;
    console.log('Connecting WS:', url);
    
    ws = new WebSocket(url);

    ws.onopen = () => {
      console.log('✅ WebSocket Connected');
      isConnected.value = true;
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        // console.log('📩 WS Msg:', msg);

        if (msg.type && listeners[msg.type]) {
          // 🚀 关键修复：把 msg.user_id 合并进 payload，否则组件不知道是谁发的消息
          const finalPayload = { ...msg.payload, user_id: msg.user_id };
          listeners[msg.type]?.forEach((cb) => cb(finalPayload));
        }
      } catch (e) {
        console.error('WS Parse Error:', e);
      }
    };

    ws.onclose = () => {
      console.log('❌ WebSocket Closed');
      isConnected.value = false;
      ws = null;
    };
    
    ws.onerror = (err) => {
      console.error('WS Error:', err);
    }
  };

  const send = (type: string, payload: any) => {
    if (ws && isConnected.value) {
      ws.send(JSON.stringify({ type, payload }));
    } else {
      console.warn('WS not connected, cannot send:', type);
    }
  };

  const on = (type: string, cb: Callback) => {
    if (!listeners[type]) listeners[type] = [];
    listeners[type].push(cb);
  };

  const off = (type: string, cb: Callback) => {
    if (!listeners[type]) return;
    listeners[type] = listeners[type].filter((c) => c !== cb);
  };

  return { isConnected, connect, send, on, off };
});