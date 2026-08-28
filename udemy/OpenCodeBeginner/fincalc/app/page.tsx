import {
  ArrowRight,
  Calculator,
  CalendarClock,
  Home as HomeIcon,
  LineChart,
  PiggyBank,
  Sparkles,
  TrendingUp,
} from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

const calculatorFeatures = [
  {
    icon: HomeIcon,
    name: "Mortgage Repayments",
    tagline: "Know exactly what a home costs you each month.",
    description:
      "Estimate monthly payments, total interest, and the full cost of a loan over its lifetime — before you commit.",
    bullets: [
      "Instant monthly payment estimates",
      "Principal vs. interest breakdown",
      "Loan term and rate comparisons",
    ],
    cta: "Estimate a mortgage",
  },
  {
    icon: TrendingUp,
    name: "Compound Interest",
    tagline: "Watch your money grow on autopilot.",
    description:
      "Project how your savings and investments compound over time with adjustable contributions, rates, and frequency.",
    bullets: [
      "Future value projections",
      "Contribution and rate sliders",
      "Year-by-year growth timeline",
    ],
    cta: "Project your growth",
  },
];

const otherFeatures = [
  {
    icon: CalendarClock,
    title: "Amortisation schedules",
    description:
      "See every payment laid out month by month so nothing is hidden.",
  },
  {
    icon: PiggyBank,
    title: "Savings goals",
    description:
      "Calculate what you need to set aside to hit a target by a date.",
  },
  {
    icon: LineChart,
    title: "Scenario comparisons",
    description:
      "A/B test rates, terms, and contributions side by side.",
  },
  {
    icon: Sparkles,
    title: "Instant, private, free",
    description:
      "All calculations run locally in your browser. No account, no data sent anywhere.",
  },
];

const heroStats = [
  { label: "2 core calculators", value: "Mortgage + Compound" },
  { label: "0 account required", value: "Works offline" },
  { label: "100% local", value: "Your data stays on your device" },
];

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col bg-background">
      {/* Header */}
      <header className="sticky top-0 z-10 w-full border-b bg-background/80 backdrop-blur">
        <div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-6">
          <a href="#" className="flex items-center gap-2">
            <span className="flex size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
              <Calculator data-icon="inline-start" />
            </span>
            <span className="font-heading text-sm font-medium">FinCalc</span>
          </a>
          <div className="flex items-center gap-3">
            <Button
              variant="ghost"
              size="sm"
              nativeButton={false}
              render={<a href="#features" />}
            >
              Features
            </Button>
            <Button size="sm">Get started</Button>
          </div>
        </div>
      </header>

      <main className="flex-1">
        {/* Hero */}
        <section className="mx-auto max-w-6xl px-6 py-20 sm:py-28">
          <div className="mx-auto flex max-w-3xl flex-col items-center gap-6 text-center">
            <Badge variant="outline" className="gap-1.5">
              <Sparkles data-icon="inline-start" />
              Financial tools that make the numbers clear
            </Badge>
            <h1 className="font-heading text-4xl font-semibold tracking-tight sm:text-5xl md:text-6xl">
              Make confident money decisions with clarity
            </h1>
            <p className="max-w-2xl text-lg text-muted-foreground">
              FinCalc brings the two formulas people actually need into one
              clean, fast tool — mortgage repayments and compound interest —
              so you can plan a house, a loan, or your savings without the
              spreadsheet headaches.
            </p>
            <div className="flex flex-wrap items-center justify-center gap-3">
              <Button size="lg">
                Try the calculators
                <ArrowRight data-icon="inline-end" />
              </Button>
              <Button
                variant="outline"
                size="lg"
                nativeButton={false}
                render={<a href="#features" />}
              >
                See how it works
              </Button>
            </div>

            <Separator className="my-8 max-w-md" />

            <dl className="grid w-full grid-cols-1 gap-6 sm:grid-cols-3">
              {heroStats.map((stat) => (
                <div
                  key={stat.label}
                  className="flex flex-col items-center gap-1"
                >
                  <dt className="text-sm font-medium text-muted-foreground">
                    {stat.label}
                  </dt>
                  <dd className="font-heading text-base">{stat.value}</dd>
                </div>
              ))}
            </dl>
          </div>
        </section>

        {/* Calculator features */}
        <section id="features" className="mx-auto max-w-6xl scroll-mt-20 px-6 pb-20">
          <div className="flex flex-col items-center gap-3 pb-12 text-center">
            <Badge variant="secondary">Core calculators</Badge>
            <h2 className="font-heading text-3xl font-semibold tracking-tight sm:text-4xl">
              Two tools, endless peace of mind
            </h2>
            <p className="max-w-2xl text-muted-foreground">
              Built around the two calculations that shape most financial
              decisions — borrowing to buy, and saving to grow.
            </p>
          </div>

          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            {calculatorFeatures.map((feature) => (
              <Card key={feature.name} className="flex flex-col">
                <CardHeader>
                  <CardAction>
                    <span className="flex size-11 items-center justify-center rounded-xl bg-primary text-primary-foreground">
                      <feature.icon />
                    </span>
                  </CardAction>
                  <Badge variant="outline">{feature.tagline}</Badge>
                  <CardTitle className="text-xl">{feature.name}</CardTitle>
                  <CardDescription className="text-base">
                    {feature.description}
                  </CardDescription>
                </CardHeader>
                <CardContent className="flex flex-1 flex-col gap-6">
                  <ul className="flex flex-col gap-3">
                    {feature.bullets.map((bullet) => (
                      <li
                        key={bullet}
                        className="flex items-start gap-3 text-sm text-muted-foreground"
                      >
                        <Badge variant="secondary" className="mt-0.5 shrink-0">
                          <CheckIcon />
                        </Badge>
                        {bullet}
                      </li>
                    ))}
                  </ul>
                </CardContent>
                <CardFooter>
                  <Button variant="outline" className="w-full">
                    {feature.cta}
                    <ArrowRight data-icon="inline-end" />
                  </Button>
                </CardFooter>
              </Card>
            ))}
          </div>
        </section>

        {/* Supporting features */}
        <section className="border-t bg-card/50">
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="flex flex-col items-center gap-3 pb-12 text-center">
              <h2 className="font-heading text-3xl font-semibold tracking-tight sm:text-4xl">
                Everything you expect, nothing you don&apos;t
              </h2>
              <p className="max-w-2xl text-muted-foreground">
                Thoughtful details that turn a calculator into an everyday
                planning tool.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
              {otherFeatures.map((feature) => (
                <Card key={feature.title} size="sm">
                  <CardHeader>
                    <CardAction>
                      <span className="flex size-9 items-center justify-center rounded-lg bg-secondary text-secondary-foreground">
                        <feature.icon data-icon="inline-start" />
                      </span>
                    </CardAction>
                    <CardTitle>{feature.title}</CardTitle>
                    <CardDescription>{feature.description}</CardDescription>
                  </CardHeader>
                </Card>
              ))}
            </div>
          </div>
        </section>

        {/* CTA */}
        <section className="mx-auto max-w-6xl px-6 py-20">
          <Card className="flex flex-col items-center overflow-hidden text-center">
            <CardHeader className="flex w-full flex-col items-center text-center">
              <Badge variant="secondary">Free forever</Badge>
              <CardTitle className="w-full max-w-3xl text-2xl font-semibold tracking-tight sm:text-3xl">
                Start planning your next big money move today
              </CardTitle>
              <CardDescription className="w-full max-w-3xl text-base">
                Estimate a mortgage or project your savings in seconds.
                No sign-up, no data sharing — just clear numbers.
              </CardDescription>
            </CardHeader>
            <CardFooter className="border-0 bg-transparent">
              <Button size="lg">
                Launch FinCalc
                <ArrowRight data-icon="inline-end" />
              </Button>
            </CardFooter>
          </Card>
        </section>
      </main>

      {/* Footer */}
      <footer className="border-t">
        <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-6 py-8 sm:flex-row">
          <p className="text-sm text-muted-foreground">
            © {new Date().getFullYear()} FinCalc. Built for people who value
            clear numbers.
          </p>
          <p className="text-sm text-muted-foreground">
            Calculators are estimates, not financial advice.
          </p>
        </div>
      </footer>
    </div>
  );
}

function CheckIcon() {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="12"
      height="12"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="3"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <path d="M20 6 9 17l-5-5" />
    </svg>
  );
}
