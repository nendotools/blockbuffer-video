export const useWebSocket = async <T>(url: string, onMessage: (data: T) => void) => {
  if (!window.WebSocket) {
    console.log("WebSocket not supported");
    return null;
  }
  const ws = new WebSocket(`http://${window.location.host}${url}`);
  ws.onmessage = (event) => {
    onMessage(JSON.parse(event.data));
  };

  return ws;
};
