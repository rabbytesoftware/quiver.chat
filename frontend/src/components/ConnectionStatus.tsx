import { Badge } from "@/components/ui/badge";

interface ConnectionStatusProps {
	isConnected: boolean;
}

export function ConnectionStatus({ isConnected }: ConnectionStatusProps) {
	return (
		<Badge
			variant={isConnected ? "default" : "destructive"}
			className={isConnected ? "bg-primary text-primary-foreground" : ""}
		>
			{isConnected ? "Connected" : "Disconnected"}
		</Badge>
	);
}
