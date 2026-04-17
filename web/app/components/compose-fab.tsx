import { useState } from "react";
import { Pencil } from "lucide-react";
import { Button } from "~/components/ui/button";
import { ComposePostDialog, type PostOptions } from "~/components/compose-post-dialog";

export function ComposeFab({
  onPost,
  currentUser,
}: {
  onPost?: (content: string, options?: PostOptions) => void;
  currentUser?: { name: string; avatarUrl?: string };
}) {
  const [open, setOpen] = useState(false);

  return (
    <>
      <Button
        onClick={() => setOpen(true)}
        className="fixed bottom-20 right-4 size-14 rounded-full shadow-lg z-20"
      >
        <Pencil className="size-5" />
      </Button>
      <ComposePostDialog
        open={open}
        onClose={() => setOpen(false)}
        onPost={onPost}
        currentUser={currentUser}
      />
    </>
  );
}
