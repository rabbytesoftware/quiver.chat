"use client";

import { useState, useEffect } from "react";
import { Moon, Sun } from "lucide-react";
import { Button } from "@/components/ui/button";

export function ThemeToggle() {
	const [theme, setTheme] = useState<"light" | "dark">("light");

	useEffect(() => {
		const savedTheme = localStorage.getItem("theme") as
			| "light"
			| "dark"
			| null;
		const systemPreference = window.matchMedia(
			"(prefers-color-scheme: dark)"
		).matches
			? "dark"
			: "light";
		const initialTheme = savedTheme || systemPreference;

		setTheme(initialTheme);
		document.documentElement.classList.toggle(
			"dark",
			initialTheme === "dark"
		);
	}, []);

	const toggleTheme = () => {
		const newTheme = theme === "light" ? "dark" : "light";
		setTheme(newTheme);
		localStorage.setItem("theme", newTheme);
		document.documentElement.classList.toggle("dark", newTheme === "dark");
	};

	return (
		<Button
			variant="outline"
			size="icon"
			onClick={toggleTheme}
			className="w-9 h-9"
			aria-label={`Switch to ${
				theme === "light" ? "dark" : "light"
			} mode`}
		>
			{theme === "light" ? (
				<Moon className="h-4 w-4" />
			) : (
				<Sun className="h-4 w-4" />
			)}
		</Button>
	);
}
