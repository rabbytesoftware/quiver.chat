import { Users } from "lucide-react";

interface OnlineCounterProps {
	count: number;
}

export function OnlineCounter({ count }: OnlineCounterProps) {
	return (
		<div className="flex items-center space-x-2 text-sm text-muted-foreground">
			<Users size={14} />
			<span className="font-mono">{count}</span>
		</div>
	);
}
