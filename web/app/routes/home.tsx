import { useState } from "react";
import { useNavigate } from "react-router";
import { ulid } from "ulid";
import { Button } from "~/components/ui/button";

export default function Home() {
  const navigate = useNavigate();
  const [roomId, setRoomId] = useState("");

  const handleCreate = () => {
    navigate(`/calls/${ulid()}`);
  };

  const handleJoin = () => {
    if (!roomId.trim()) return;
    navigate(`/calls/${roomId.trim()}`);
  };

  return (
    <div className="w-full min-h-full flex flex-col items-center justify-center px-6">
      <div className="w-full max-w-sm flex flex-col gap-6">
        <Button onClick={handleCreate} className="w-full h-11">
          新しい通話を作成
        </Button>

        <div className="flex flex-col gap-2">
          <input
            className="w-full h-11 px-3 rounded-md border border-input bg-background text-sm"
            placeholder="Room ID で参加"
            value={roomId}
            onChange={(e) => setRoomId(e.target.value)}
          />
          <Button
            variant="outline"
            onClick={handleJoin}
            disabled={!roomId.trim()}
            className="w-full h-11"
          >
            参加
          </Button>
        </div>
      </div>
    </div>
  );
}
