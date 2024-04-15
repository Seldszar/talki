import { useEffect, useState } from "preact/hooks";

import { Message } from "../types/socket";
import { Speaker } from "../types/speaker";

export function useSpeakers() {
  const [speakers, setSpeakers] = useState(new Array<Speaker>());

  useEffect(() => {
    const url = new URL("/socket", location.href.replace("http", "ws"));

    let socket: WebSocket | undefined;
    let timeoutId: number | undefined;

    const connect = () => {
      socket?.close(1_000);

      if (timeoutId) {
        clearTimeout(timeoutId);
      }

      socket = new WebSocket(url);

      socket.addEventListener("close", (event) => {
        if (event.code === 1_000) {
          return;
        }

        timeoutId = setTimeout(connect, 1_000);
      });

      socket.addEventListener("message", (event) => {
        const data = JSON.parse(event.data) as Message;

        if (data.speakers) {
          setSpeakers(data.speakers.sort((a, b) => a.displayName.localeCompare(b.displayName)));
        }
      });
    };

    connect();

    return () => {
      socket?.close(1_000);

      if (timeoutId) {
        clearTimeout(timeoutId);
      }
    };
  }, []);

  return speakers;
}
