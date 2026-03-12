import { Link } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { ArrowRight, Phone, Video } from "lucide-react";

export interface CallParticipant {
  name: string;
  avatarUrl?: string;
}

export interface Call {
  id: string;
  name: string;
  type: "audio" | "video";
  participants: CallParticipant[];
}

export function CallList({ calls }: { calls: Call[] }) {
  return (
    <div className="flex items-center gap-4 px-4 py-5 overflow-x-auto border-b">
      {calls.map((call) => (
        <CallItem key={call.id} call={call} />
      ))}
      <Link
        to="/calls"
        className="flex flex-col items-center gap-1.5 shrink-0"
      >
        <div className="size-14 rounded-full bg-muted flex items-center justify-center text-muted-foreground hover:text-foreground transition-colors">
          <ArrowRight className="size-5" />
        </div>
        <span className="text-xs text-muted-foreground whitespace-nowrap">
          もっと見る
        </span>
      </Link>
    </div>
  );
}

function CallItem({ call }: { call: Call }) {
  const TypeIcon = call.type === "video" ? Video : Phone;

  return (
    <button className="flex flex-col items-center gap-1.5 shrink-0">
      <div className="relative">
        <GroupAvatar participants={call.participants} />
        <div className="absolute -bottom-0.5 -right-0.5 size-5 rounded-full bg-primary flex items-center justify-center">
          <TypeIcon className="size-3 text-primary-foreground" />
        </div>
      </div>
      <span className="text-xs truncate max-w-16">{call.name}</span>
    </button>
  );
}

function GroupAvatar({
  participants,
}: {
  participants: CallParticipant[];
}) {
  const count = Math.min(participants.length, 4);

  if (count <= 1) {
    const p = participants[0];
    return (
      <Avatar className="size-14 ring-2 ring-primary ring-offset-2 ring-offset-background">
        <AvatarImage src={p?.avatarUrl} alt={p?.name} />
        <AvatarFallback>{p?.name.slice(0, 2)}</AvatarFallback>
      </Avatar>
    );
  }

  return (
    <div className="size-14 rounded-full ring-2 ring-primary ring-offset-2 ring-offset-background overflow-hidden relative bg-muted">
      {count === 2 && (
        <>
          <MiniAvatar participant={participants[0]} className="absolute top-0 left-0 w-1/2 h-full rounded-none" />
          <MiniAvatar participant={participants[1]} className="absolute top-0 right-0 w-1/2 h-full rounded-none" />
          <div className="absolute inset-y-0 left-1/2 w-px bg-background" />
        </>
      )}
      {count === 3 && (
        <>
          <MiniAvatar participant={participants[0]} className="absolute top-0 left-0 w-1/2 h-full rounded-none" />
          <MiniAvatar participant={participants[1]} className="absolute top-0 right-0 w-1/2 h-1/2 rounded-none" />
          <MiniAvatar participant={participants[2]} className="absolute bottom-0 right-0 w-1/2 h-1/2 rounded-none" />
          <div className="absolute inset-y-0 left-1/2 w-px bg-background" />
          <div className="absolute top-1/2 left-1/2 right-0 h-px bg-background" />
        </>
      )}
      {count === 4 && (
        <>
          <MiniAvatar participant={participants[0]} className="absolute top-0 left-0 w-1/2 h-1/2 rounded-none" />
          <MiniAvatar participant={participants[1]} className="absolute top-0 right-0 w-1/2 h-1/2 rounded-none" />
          <MiniAvatar participant={participants[2]} className="absolute bottom-0 left-0 w-1/2 h-1/2 rounded-none" />
          <MiniAvatar participant={participants[3]} className="absolute bottom-0 right-0 w-1/2 h-1/2 rounded-none" />
          <div className="absolute inset-y-0 left-1/2 w-px bg-background" />
          <div className="absolute inset-x-0 top-1/2 h-px bg-background" />
        </>
      )}
    </div>
  );
}

function MiniAvatar({
  participant,
  className,
}: {
  participant: CallParticipant;
  className?: string;
}) {
  return (
    <div className={`overflow-hidden flex items-center justify-center bg-muted ${className ?? ""}`}>
      {participant.avatarUrl ? (
        <img
          src={participant.avatarUrl}
          alt={participant.name}
          className="w-full h-full object-cover"
        />
      ) : (
        <span className="text-[10px] text-muted-foreground">
          {participant.name.slice(0, 1)}
        </span>
      )}
    </div>
  );
}
