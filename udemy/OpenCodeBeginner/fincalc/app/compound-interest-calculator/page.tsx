"use client"

import { useMemo, useState } from "react"
import Link from "next/link"
import {
  Area,
  AreaChart,
  CartesianGrid,
  Pie,
  PieChart,
  XAxis,
  YAxis,
} from "recharts"
import {
  Calculator,
  FileSpreadsheet,
  LineChart as LineChartIcon,
  TrendingUp,
} from "lucide-react"

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart"
import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from "@/components/ui/empty"
import {
  Field,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group"

const chartConfig = {
  contributions: {
    label: "Contributions",
    color: "var(--chart-1)",
  },
  interest: {
    label: "Interest",
    color: "var(--chart-2)",
  },
  balance: {
    label: "Balance",
    color: "var(--chart-3)",
  },
} satisfies ChartConfig

const currencyFormatter = new Intl.NumberFormat("en-GB", {
  style: "currency",
  currency: "USD",
  currencyDisplay: "narrowSymbol",
  maximumFractionDigits: 0,
})

const currencyFormatterFull = new Intl.NumberFormat("en-GB", {
  style: "currency",
  currency: "USD",
  currencyDisplay: "narrowSymbol",
  maximumFractionDigits: 2,
})

const compoundFrequencyOptions = [
  { value: "1", label: "Annually" },
  { value: "4", label: "Quarterly" },
  { value: "12", label: "Monthly" },
]

const termOptions = ["5", "10", "15", "20", "25", "30"]

type YearRow = {
  year: number
  openingBalance: number
  contributions: number
  interestEarned: number
  closingBalance: number
}

type YearGroup = {
  year: number
  rows: YearRow[]
}

type GrowthPoint = {
  year: number
  balance: number
}

export default function CompoundInterestCalculator() {
  const [initialInvestment, setInitialInvestment] = useState("10000")
  const [monthlyContribution, setMonthlyContribution] = useState("200")
  const [interestRate, setInterestRate] = useState("7.0")
  const [compoundFrequency, setCompoundFrequency] = useState(["12"])
  const [term, setTerm] = useState(["10"])

  const results = useMemo(() => {
    const initial = parseFloat(initialInvestment) || 0
    const monthly = parseFloat(monthlyContribution) || 0
    const rate = parseFloat(interestRate) || 0
    const years = parseInt(term[0] || "10", 10)

    const r = rate / 100
    const totalMonths = years * 12

    if (initial === 0 && monthly === 0) {
      return {
        futureValue: 0,
        totalContributions: 0,
        totalInterest: 0,
        chartData: [] as { name: string; value: number; fill: string }[],
        growthData: [] as GrowthPoint[],
        schedule: [] as YearRow[],
        yearGroups: [] as YearGroup[],
      }
    }

    const monthlyRate = r / 12
    let balance = initial
    let totalContributions = initial
    let totalInterest = 0

    const schedule: YearRow[] = []
    const growthData: GrowthPoint[] = [{ year: 0, balance: initial }]

    for (let month = 1; month <= totalMonths; month++) {
      balance += monthly
      totalContributions += monthly

      const monthInterest = balance * monthlyRate
      balance += monthInterest
      totalInterest += monthInterest

      if (month % 12 === 0) {
        const year = month / 12
        const prevBalance = balance - monthInterest - monthly
        schedule.push({
          year,
          openingBalance: prevBalance,
          contributions: monthly * 12,
          interestEarned: monthInterest * 12,
          closingBalance: balance,
        })
        growthData.push({ year, balance })
      }
    }

    const futureValue = balance

    const chartData = [
      {
        name: "Contributions",
        value: Math.round(totalContributions),
        fill: "var(--color-contributions)",
      },
      {
        name: "Interest",
        value: Math.round(totalInterest),
        fill: "var(--color-interest)",
      },
    ]

    const yearGroups = schedule.reduce<YearGroup[]>((groups, row) => {
      const existing = groups[groups.length - 1]
      if (existing && existing.year === row.year) {
        existing.rows.push(row)
      } else {
        groups.push({ year: row.year, rows: [row] })
      }
      return groups
    }, [])

    return {
      futureValue,
      totalContributions,
      totalInterest,
      chartData,
      growthData,
      schedule,
      yearGroups,
    }
  }, [initialInvestment, monthlyContribution, interestRate, term])

  return (
    <div className="flex min-h-screen flex-col bg-background">
      {/* Header */}
      <header className="sticky top-0 z-10 w-full border-b bg-background/80 backdrop-blur">
        <div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-6">
          <Link href="/" className="flex items-center gap-2">
            <span className="flex size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
              <Calculator data-icon="inline-start" />
            </span>
            <span className="font-heading text-sm font-medium">FinCalc</span>
          </Link>
          <Button size="sm" nativeButton={false} render={<Link href="/" />}>
            All calculators
          </Button>
        </div>
      </header>

      <main className="flex-1">
        <div className="mx-auto max-w-6xl px-6 py-12">
          {/* Page header */}
          <div className="flex flex-col gap-4 pb-10">
            <Badge variant="secondary" className="w-fit gap-1.5">
              <TrendingUp data-icon="inline-start" />
              Compound interest calculator
            </Badge>
            <h1 className="font-heading text-3xl font-semibold tracking-tight sm:text-4xl">
              Compound Interest Calculator
            </h1>
            <p className="max-w-2xl text-muted-foreground">
              See how your savings and investments grow over time with the power
              of compound interest — and plan your path to a target balance.
            </p>
          </div>

          {/* Tabs */}
          <Tabs defaultValue="results">
            <TabsList>
              <TabsTrigger value="results">Results</TabsTrigger>
              <TabsTrigger value="growth-chart">
                <LineChartIcon data-icon="inline-start" />
                Growth chart
              </TabsTrigger>
              <TabsTrigger value="growth">
                <FileSpreadsheet data-icon="inline-start" />
                Year-by-year breakdown
              </TabsTrigger>
            </TabsList>

            <TabsContent value="results">
              {/* Calculator grid */}
              <div className="grid grid-cols-1 gap-6 lg:grid-cols-5">
                {/* Input form */}
                <Card className="lg:col-span-3">
                  <CardHeader>
                    <CardTitle>Investment details</CardTitle>
                    <CardDescription>
                      Enter your savings details to project compound growth.
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <FieldGroup>
                      <Field>
                        <FieldLabel htmlFor="initial-investment">
                          Initial investment
                        </FieldLabel>
                        <div className="relative">
                          <span className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-sm text-muted-foreground">
                            $
                          </span>
                          <Input
                            id="initial-investment"
                            type="number"
                            className="pl-6"
                            value={initialInvestment}
                            onChange={(e) =>
                              setInitialInvestment(e.target.value)
                            }
                            min={0}
                            step={1000}
                          />
                        </div>
                      </Field>

                      <Field>
                        <FieldLabel htmlFor="monthly-contribution">
                          Monthly contribution
                        </FieldLabel>
                        <div className="relative">
                          <span className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-sm text-muted-foreground">
                            $
                          </span>
                          <Input
                            id="monthly-contribution"
                            type="number"
                            className="pl-6"
                            value={monthlyContribution}
                            onChange={(e) =>
                              setMonthlyContribution(e.target.value)
                            }
                            min={0}
                            step={50}
                          />
                        </div>
                      </Field>

                      <Field>
                        <FieldLabel htmlFor="interest-rate">
                          Annual interest rate
                        </FieldLabel>
                        <div className="relative">
                          <Input
                            id="interest-rate"
                            type="number"
                            className="pr-7"
                            value={interestRate}
                            onChange={(e) => setInterestRate(e.target.value)}
                            min={0}
                            max={30}
                            step={0.1}
                          />
                          <span className="pointer-events-none absolute right-2.5 top-1/2 -translate-y-1/2 text-sm text-muted-foreground">
                            %
                          </span>
                        </div>
                      </Field>

                      <Field>
                        <FieldLabel>Compound frequency</FieldLabel>
                        <ToggleGroup
                          defaultValue={["12"]}
                          value={compoundFrequency}
                          onValueChange={(v) => {
                            if (v.length > 0) setCompoundFrequency(v)
                          }}
                          spacing={2}
                        >
                          {compoundFrequencyOptions.map((opt) => (
                            <ToggleGroupItem key={opt.value} value={opt.value}>
                              {opt.label}
                            </ToggleGroupItem>
                          ))}
                        </ToggleGroup>
                      </Field>

                      <Field>
                        <FieldLabel>Time period</FieldLabel>
                        <ToggleGroup
                          defaultValue={["10"]}
                          value={term}
                          onValueChange={(v) => {
                            if (v.length > 0) setTerm(v)
                          }}
                          spacing={2}
                        >
                          {termOptions.map((t) => (
                            <ToggleGroupItem key={t} value={t}>
                              {t} yr
                            </ToggleGroupItem>
                          ))}
                        </ToggleGroup>
                      </Field>
                    </FieldGroup>
                  </CardContent>
                </Card>

                {/* Results */}
                <div className="flex flex-col gap-6 lg:col-span-2">
                  {/* Future value */}
                  <Card>
                    <CardHeader>
                      <CardDescription>Future balance</CardDescription>
                      <CardAction>
                        <CardTitle className="text-3xl font-semibold tabular-nums">
                          {currencyFormatterFull.format(results.futureValue)}
                        </CardTitle>
                      </CardAction>
                    </CardHeader>
                  </Card>

                  {/* Totals */}
                  <Card>
                    <CardHeader>
                      <CardTitle className="text-base">
                        Growth breakdown
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="flex flex-col gap-4">
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-muted-foreground">
                          Total contributions
                        </span>
                        <span className="font-medium tabular-nums">
                          {currencyFormatter.format(
                            results.totalContributions
                          )}
                        </span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-muted-foreground">
                          Total interest earned
                        </span>
                        <span className="font-medium tabular-nums text-emerald-600">
                          {currencyFormatter.format(results.totalInterest)}
                        </span>
                      </div>
                      <div className="border-t pt-4">
                        <div className="flex items-center justify-between">
                          <span className="text-sm font-medium">
                            Final balance
                          </span>
                          <span className="font-heading text-lg font-semibold tabular-nums">
                            {currencyFormatter.format(results.futureValue)}
                          </span>
                        </div>
                      </div>
                    </CardContent>
                  </Card>

                  {/* Pie chart */}
                  {results.chartData.length > 0 &&
                    results.totalInterest > 0 && (
                      <Card>
                        <CardHeader>
                          <CardTitle className="text-base">
                            Contributions vs interest
                          </CardTitle>
                        </CardHeader>
                        <CardContent>
                          <ChartContainer
                            config={chartConfig}
                            className="mx-auto aspect-square max-h-[250px]"
                          >
                            <PieChart accessibilityLayer>
                              <ChartTooltip
                                content={
                                  <ChartTooltipContent
                                    nameKey="name"
                                    formatter={(value) =>
                                      currencyFormatter.format(
                                        value as number
                                      )
                                    }
                                  />
                                }
                              />
                              <Pie
                                data={results.chartData}
                                dataKey="value"
                                nameKey="name"
                                innerRadius={60}
                                strokeWidth={5}
                              />
                              <ChartLegend
                                content={
                                  <ChartLegendContent nameKey="name" />
                                }
                              />
                            </PieChart>
                          </ChartContainer>
                        </CardContent>
                      </Card>
                    )}
                </div>
              </div>
            </TabsContent>

            <TabsContent value="growth-chart">
              {results.growthData.length > 0 ? (
                <Card>
                  <CardHeader>
                    <CardTitle className="text-base">
                      Investment growth
                    </CardTitle>
                    <CardDescription>
                      Projected balance year by year over your {term[0]}-year
                      investment at {interestRate || "0"}% interest.
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <ChartContainer
                      config={chartConfig}
                      className="h-[320px] w-full"
                    >
                      <AreaChart
                        accessibilityLayer
                        data={results.growthData}
                        margin={{ left: 8, right: 8 }}
                      >
                        <CartesianGrid vertical={false} />
                        <XAxis
                          dataKey="year"
                          tickLine={false}
                          axisLine={false}
                          tickMargin={8}
                          tickFormatter={(value) => `Year ${value}`}
                        />
                        <YAxis
                          tickLine={false}
                          axisLine={false}
                          width={60}
                          tickFormatter={(value) =>
                            currencyFormatter.format(value as number)
                          }
                        />
                        <ChartTooltip
                          content={
                            <ChartTooltipContent
                              formatter={(value) =>
                                currencyFormatterFull.format(value as number)
                              }
                            />
                          }
                        />
                        <Area
                          dataKey="balance"
                          type="monotone"
                          fill="var(--color-balance)"
                          fillOpacity={0.3}
                          stroke="var(--color-balance)"
                          strokeWidth={2}
                        />
                      </AreaChart>
                    </ChartContainer>
                  </CardContent>
                </Card>
              ) : (
                <Card>
                  <CardContent>
                    <Empty>
                      <EmptyHeader>
                        <EmptyMedia variant="icon">
                          <LineChartIcon />
                        </EmptyMedia>
                        <EmptyTitle>No growth chart available</EmptyTitle>
                        <EmptyDescription>
                          Enter an initial investment or monthly contribution
                          to see how your money grows over time.
                        </EmptyDescription>
                      </EmptyHeader>
                    </Empty>
                  </CardContent>
                </Card>
              )}
            </TabsContent>

            <TabsContent value="growth">
              {results.schedule.length > 0 ? (
                <Card>
                  <CardHeader>
                    <CardTitle className="text-base">
                      Year-by-year growth
                    </CardTitle>
                    <CardDescription>
                      Annual breakdown of your {term[0]}-year investment at{" "}
                      {interestRate || "0"}% compounded{" "}
                      {compoundFrequencyOptions.find(
                        (o) => o.value === compoundFrequency[0]
                      )?.label.toLowerCase() || "monthly"}
                      .
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="flex flex-col gap-3">
                      <p className="text-sm text-muted-foreground">
                        Expand a year to see the opening balance, contributions
                        added, interest earned, and closing balance for that
                        year.
                      </p>
                      <Accordion
                        multiple
                        defaultValue={["year-1"]}
                        className="rounded-md border"
                      >
                        {results.yearGroups.map((group) => (
                          <AccordionItem
                            key={group.year}
                            value={`year-${group.year}`}
                          >
                            <AccordionTrigger>
                              <span className="flex w-full items-center justify-between gap-2">
                                <span className="font-heading font-semibold uppercase tracking-wide text-muted-foreground">
                                  Year {group.year}
                                </span>
                                <span className="text-xs font-normal normal-case tracking-normal text-muted-foreground/80">
                                  {currencyFormatter.format(
                                    group.rows[0]?.closingBalance ?? 0
                                  )}
                                </span>
                              </span>
                            </AccordionTrigger>
                            <AccordionContent>
                              <Table>
                                <TableHeader>
                                  <TableRow>
                                    <TableHead>Opening balance</TableHead>
                                    <TableHead className="text-right">
                                      Contributions
                                    </TableHead>
                                    <TableHead className="text-right">
                                      Interest
                                    </TableHead>
                                    <TableHead className="text-right">
                                      Closing balance
                                    </TableHead>
                                  </TableRow>
                                </TableHeader>
                                <TableBody>
                                  {group.rows.map((row) => (
                                    <TableRow key={row.year}>
                                      <TableCell className="font-medium tabular-nums">
                                        {currencyFormatterFull.format(
                                          row.openingBalance
                                        )}
                                      </TableCell>
                                      <TableCell className="text-right tabular-nums">
                                        {currencyFormatterFull.format(
                                          row.contributions
                                        )}
                                      </TableCell>
                                      <TableCell className="text-right tabular-nums">
                                        {currencyFormatterFull.format(
                                          row.interestEarned
                                        )}
                                      </TableCell>
                                      <TableCell className="text-right tabular-nums">
                                        {currencyFormatterFull.format(
                                          row.closingBalance
                                        )}
                                      </TableCell>
                                    </TableRow>
                                  ))}
                                </TableBody>
                              </Table>
                            </AccordionContent>
                          </AccordionItem>
                        ))}
                      </Accordion>
                    </div>
                  </CardContent>
                </Card>
              ) : (
                <Card>
                  <CardContent>
                    <Empty>
                      <EmptyHeader>
                        <EmptyMedia variant="icon">
                          <FileSpreadsheet />
                        </EmptyMedia>
                        <EmptyTitle>
                          No growth breakdown available
                        </EmptyTitle>
                        <EmptyDescription>
                          Enter an initial investment or monthly contribution
                          to see how your money grows year by year.
                        </EmptyDescription>
                      </EmptyHeader>
                    </Empty>
                  </CardContent>
                </Card>
              )}
            </TabsContent>
          </Tabs>
        </div>
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
  )
}
