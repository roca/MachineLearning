import Link from "next/link";

const calculators = [
  {
    href: "/mortgage",
    title: "Mortgage Repayment Calculator",
    description:
      "Work out your monthly repayments, total interest, and the true cost of a home loan in seconds.",
    points: [
      "Monthly repayment breakdown",
      "Total interest over the loan term",
      "Compare different rates and terms",
    ],
    icon: (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-6 w-6">
        <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 12l8.954-8.955a1.5 1.5 0 012.09 0l8.455 8.455" />
        <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 10.5V21h15V10.5" />
        <path strokeLinecap="round" strokeLinejoin="round" d="M9.75 21v-6h4.5v6" />
      </svg>
    ),
  },
  {
    href: "/compound-interest",
    title: "Compound Interest Calculator",
    description:
      "See how your savings and investments grow over time when compound interest works in your favor.",
    points: [
      "Project growth year by year",
      "Adjust principal, rate, and duration",
      "Visualize the power of compounding",
    ],
    icon: (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-6 w-6">
        <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 18L9 11.25l4.306 4.307a11.95 11.95 0 015.814-5.519l2.74-1.22m0 0l-5.94-2.28m5.94 2.28l-2.28 5.941" />
      </svg>
    ),
  },
];

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col bg-zinc-50 font-sans dark:bg-black">
      <header className="sticky top-0 z-10 border-b border-zinc-200 bg-white/80 backdrop-blur dark:border-zinc-800 dark:bg-black/80">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
          <span className="text-lg font-semibold tracking-tight text-zinc-950 dark:text-zinc-50">
            FinCalc
          </span>
          <nav className="flex items-center gap-6 text-sm font-medium text-zinc-600 dark:text-zinc-400">
            <Link href="/mortgage" className="transition-colors hover:text-zinc-950 dark:hover:text-zinc-50">
              Mortgage
            </Link>
            <Link href="/compound-interest" className="transition-colors hover:text-zinc-950 dark:hover:text-zinc-50">
              Compound Interest
            </Link>
          </nav>
        </div>
      </header>

      <main className="flex-1">
        <section className="mx-auto flex max-w-6xl flex-col items-center px-6 py-24 text-center sm:py-32">
          <p className="rounded-full border border-emerald-600/20 bg-emerald-600/10 px-3 py-1 text-xs font-semibold tracking-wide text-emerald-700 uppercase dark:border-emerald-500/20 dark:bg-emerald-500/10 dark:text-emerald-400">
            Free · Fast · No sign-up
          </p>
          <h1 className="mt-6 max-w-2xl text-4xl font-semibold leading-tight tracking-tight text-zinc-950 sm:text-5xl dark:text-zinc-50">
            Make smarter money decisions with simple financial tools
          </h1>
          <p className="mt-6 max-w-xl text-lg leading-8 text-zinc-600 dark:text-zinc-400">
            Plan a home loan or grow your savings. FinCalc gives you clear,
            instant answers without spreadsheets or jargon.
          </p>
          <div className="mt-8 flex flex-col gap-3 sm:flex-row">
            <Link
              href="#calculators"
              className="flex h-12 items-center justify-center rounded-full bg-zinc-950 px-6 font-medium text-white transition-colors hover:bg-zinc-700 dark:bg-zinc-50 dark:text-zinc-950 dark:hover:bg-zinc-300"
            >
              Explore Calculators
            </Link>
          </div>
        </section>

        <section id="calculators" className="mx-auto max-w-6xl px-6 pb-24 scroll-mt-16">
          <div className="grid gap-6 md:grid-cols-2">
            {calculators.map((calc) => (
              <Link
                key={calc.href}
                href={calc.href}
                className="group flex flex-col rounded-2xl border border-zinc-200 bg-white p-8 transition-all hover:-translate-y-1 hover:border-zinc-300 hover:shadow-lg dark:border-zinc-800 dark:bg-zinc-900 dark:hover:border-zinc-700"
              >
                <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-emerald-600/10 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400">
                  {calc.icon}
                </div>
                <h2 className="mt-6 text-xl font-semibold tracking-tight text-zinc-950 dark:text-zinc-50">
                  {calc.title}
                </h2>
                <p className="mt-3 leading-7 text-zinc-600 dark:text-zinc-400">
                  {calc.description}
                </p>
                <ul className="mt-6 space-y-2.5">
                  {calc.points.map((point) => (
                    <li key={point} className="flex items-start gap-2.5 text-sm text-zinc-600 dark:text-zinc-400">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="mt-0.5 h-4 w-4 shrink-0 text-emerald-600 dark:text-emerald-400">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                      </svg>
                      {point}
                    </li>
                  ))}
                </ul>
                <span className="mt-8 inline-flex items-center gap-1.5 text-sm font-medium text-zinc-950 transition-colors group-hover:text-emerald-700 dark:text-zinc-50 dark:group-hover:text-emerald-400">
                  Open calculator
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="h-4 w-4 transition-transform group-hover:translate-x-0.5">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
                  </svg>
                </span>
              </Link>
            ))}
          </div>
        </section>

        <section className="border-t border-zinc-200 bg-white dark:border-zinc-800 dark:bg-zinc-900">
          <div className="mx-auto grid max-w-6xl gap-8 px-6 py-16 sm:grid-cols-3">
            {[
              {
                title: "Instant results",
                description:
                  "Every calculation updates in real time as you type. No page reloads, no waiting.",
              },
              {
                title: "Accurate math",
                description:
                  "Standard amortization and compound interest formulas you can trust for planning.",
              },
              {
                title: "Private by design",
                description:
                  "Everything runs in your browser. Your numbers never leave your device.",
              },
            ].map((item) => (
              <div key={item.title}>
                <h3 className="font-semibold text-zinc-950 dark:text-zinc-50">{item.title}</h3>
                <p className="mt-2 text-sm leading-6 text-zinc-600 dark:text-zinc-400">
                  {item.description}
                </p>
              </div>
            ))}
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
