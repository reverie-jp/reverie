import { Mic, MicOff, type LucideIcon } from "lucide-react";

export function CallBar({
  title,
  subtitle,
  micOn,
  onToggleMic,
  onTapCenter,
  onEnd,
  EndIcon,
}: {
  title: string;
  subtitle: string;
  micOn: boolean;
  onToggleMic: () => void;
  onTapCenter: () => void;
  onEnd: () => void;
  EndIcon: LucideIcon;
}) {
  return (
    <div className="w-full shrink-0">
      <div className="bg-muted border-b px-3 py-2.5 rounded-b-2xl shadow-md">
        <div className="flex items-center gap-2">
          <button
            className={`size-9 rounded-full flex items-center justify-center transition-colors ${
              micOn
                ? "bg-background hover:bg-background/80 text-foreground"
                : "bg-red-500 hover:bg-red-600 text-white"
            }`}
            onClick={(e) => {
              e.stopPropagation();
              onToggleMic();
            }}
          >
            {micOn ? <Mic className="size-4" /> : <MicOff className="size-4" />}
          </button>
          <button
            className="flex-1 flex flex-col items-center min-w-0 py-0.5"
            onClick={onTapCenter}
          >
            <span className="text-sm font-bold truncate max-w-full text-foreground">
              {title}
            </span>
            <span className="text-xs text-muted-foreground">{subtitle}</span>
          </button>
          <button
            className="size-9 rounded-full bg-red-500 hover:bg-red-600 text-white flex items-center justify-center transition-colors"
            onClick={(e) => {
              e.stopPropagation();
              onEnd();
            }}
          >
            <EndIcon className="size-4" />
          </button>
        </div>
      </div>
    </div>
  );
}
