"use client"

import Link from "next/link"
import { Calculator } from "lucide-react"
import { ThemeToggle } from "@/components/theme-toggle"

export function Header({ rightContent }: { rightContent?: React.ReactNode }) {
  return (
    <header className="sticky top-0 z-10 w-full border-b bg-background/80 backdrop-blur">
      <div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-6">
        <Link href="/" className="flex items-center gap-2">
          <span className="flex size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <Calculator data-icon="inline-start" />
          </span>
          <span className="font-heading text-sm font-medium">FinCalc</span>
        </Link>
        <div className="flex items-center gap-3">
          {rightContent}
          <ThemeToggle />
        </div>
      </div>
    </header>
  )
}
