import Link from "next/link";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Test — FinCalc",
  description: "Example card page for testing.",
};

export default function Test() {
  return (
    <div className="flex min-h-screen flex-col bg-zinc-50 font-sans dark:bg-black">
      <header className="sticky top-0 z-10 border-b border-zinc-200 bg-white/80 backdrop-blur dark:border-zinc-800 dark:bg-black/80">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
          <span className="text-lg font-semibold tracking-tight text-zinc-950 dark:text-zinc-50">
            FinCalc
          </span>
          <nav className="flex items-center gap-6 text-sm font-medium text-zinc-600 dark:text-zinc-400">
            <Link href="/" className="transition-colors hover:text-zinc-950 dark:hover:text-zinc-50">
              Home
            </Link>
          </nav>
        </div>
      </header>

      <main className="flex-1">
        <section className="mx-auto max-w-6xl px-6 py-24">
          <h1 className="text-3xl font-semibold tracking-tight text-zinc-950 sm:text-4xl dark:text-zinc-50">
            Test Page
          </h1>
          <p className="mt-4 max-w-xl text-lg leading-8 text-zinc-600 dark:text-zinc-400">
            An example card to verify layout and styling.
          </p>

          <div className="mt-12 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            <div className="group flex flex-col rounded-2xl border border-pink-200 bg-pink-50 p-8 transition-all hover:-translate-y-1 hover:border-pink-300 hover:shadow-lg dark:border-pink-900 dark:bg-pink-950 dark:hover:border-pink-800">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-pink-600/10 text-pink-700 dark:bg-pink-500/10 dark:text-pink-400">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-6 w-6">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09z" />
                </svg>
              </div>
              <h2 className="mt-6 text-xl font-semibold tracking-tight text-zinc-950 dark:text-zinc-50">
                Example Card
              </h2>
              <p className="mt-3 leading-7 text-zinc-600 dark:text-zinc-400">
                This is a sample card component demonstrating the standard card
                styling used across FinCalc.
              </p>
              <ul className="mt-6 space-y-2.5">
                {["Rounded corners", "Hover elevation", "Dark mode ready"].map((point) => (
                  <li key={point} className="flex items-start gap-2.5 text-sm text-zinc-600 dark:text-zinc-400">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="mt-0.5 h-4 w-4 shrink-0 text-pink-600 dark:text-pink-400">
                      <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                    </svg>
                    {point}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </section>
      </main>

      <footer className="border-t border-zinc-200 dark:border-zinc-800">
        <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-6 py-8 text-sm text-zinc-600 sm:flex-row dark:text-zinc-400">
          <span>FinCalc</span>
          <span>Estimates are for guidance only, not financial advice.</span>
        </div>
      </footer>
    </div>
  );
}
