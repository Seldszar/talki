import { cx } from "classix";
import { useMemo } from "preact/hooks";

import { useSpeakers } from "../hooks/speakers";

export function App() {
  const speakers = useSpeakers();

  const searchParams = useMemo(
    () => new URLSearchParams(location.search),
    [location.search],
  );

  return (
    <div class="Voice_voiceContainer__adk9M voice_container">
      <ul class="Voice_voiceStates__a121W voice_states">
        {speakers.map((speaker) => (
          <li key={speaker.id} class={cx("Voice_voiceState__OCoZh voice_state", speaker.speaking && "wrapper_speaking")} data-userid={speaker.id}>
            <img class={cx("Voice_avatar__htiqH voice_avatar", speaker.speaking && "Voice_avatarSpeaking__lE+4m")} src={speaker.avatarUrl} />
            <div class="Voice_user__8fGwX voice_username">
              <span class="Voice_name__TALd9" style={{ backgroundColor: searchParams.get("bg_color"), color: searchParams.get("text_color"), fontSize: Number.parseInt(searchParams.get("text_size") ?? "") }}>{speaker.displayName}</span>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
