"use client"

import { useMemo, useState } from "react"
import Link from "next/link"
import { Pie, PieChart } from "recharts"
import { Calculator, FileSpreadsheet, Home as HomeIcon } from "lucide-react"

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
  FieldDescription,
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
  principal: {
    label: "Principal",
    color: "var(--chart-1)",
  },
  interest: {
    label: "Interest",
    color: "var(--chart-2)",
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

const loanTermOptions = ["10", "15", "20", "25", "30"]

type ScheduleRow = {
  month: number
  payment: number
  interest: number
  principal: number
  balance: number
}

type YearGroup = {
  year: number
  firstMonth: number
  lastMonth: number
  rows: ScheduleRow[]
}

export default function MortgageRepaymentsCalculator() {
  const [homePrice, setHomePrice] = useState("250000")
  const [deposit, setDeposit] = useState("50000")
  const [interestRate, setInterestRate] = useState("5.0")
  const [loanTerm, setLoanTerm] = useState(["25"])
  const [repaymentType, setRepaymentType] = useState(["repayment"])

  const results = useMemo(() => {
    const price = parseFloat(homePrice) || 0
    const dep = parseFloat(deposit) || 0
    const rate = parseFloat(interestRate) || 0
    const termYears = parseInt(loanTerm[0] || "25", 10)
    const isRepayment = repaymentType[0] === "repayment"

    const principal = Math.max(price - dep, 0)
    const monthlyRate = rate / 100 / 12
    const totalMonths = termYears * 12

    if (principal === 0 || monthlyRate === 0) {
      return {
        monthlyPayment: 0,
        totalPaid: 0,
        totalInterest: 0,
        principal,
        chartData: [] as { name: string; value: number; fill: string }[],
        schedule: [] as ScheduleRow[],
        yearGroups: [] as YearGroup[],
      }
    }

    let monthlyPayment: number

    if (isRepayment) {
      monthlyPayment =
        (principal * (monthlyRate * Math.pow(1 + monthlyRate, totalMonths))) /
        (Math.pow(1 + monthlyRate, totalMonths) - 1)
    } else {
      monthlyPayment = principal * monthlyRate
    }

    const totalPaid = monthlyPayment * totalMonths
    const totalInterest = isRepayment ? totalPaid - principal : totalPaid

    const chartData = [
      {
        name: "Principal",
        value: Math.round(principal),
        fill: "var(--color-principal)",
      },
      {
        name: "Interest",
        value: Math.round(totalInterest),
        fill: "var(--color-interest)",
      },
    ]

    const schedule: ScheduleRow[] = []

    if (isRepayment) {
      let balance = principal
      for (let month = 1; month <= totalMonths; month++) {
        const interest = balance * monthlyRate
        const principalPaid = monthlyPayment - interest
        balance = Math.max(balance - principalPaid, 0)
        schedule.push({
          month,
          payment: monthlyPayment,
          interest,
          principal: principalPaid,
          balance,
        })
      }
    }

    const yearGroups = schedule.reduce<YearGroup[]>((groups, row) => {
      const year = Math.ceil(row.month / 12)
      const existing = groups[groups.length - 1]
      if (existing && existing.year === year) {
        existing.rows.push(row)
      } else {
        groups.push({
          year,
          firstMonth: row.month,
          lastMonth: row.month,
          rows: [row],
        })
      }
      return groups
    }, [])

    return {
      monthlyPayment,
      totalPaid,
      totalInterest,
      principal,
      chartData,
      schedule,
      yearGroups,
    }
  }, [homePrice, deposit, interestRate, loanTerm, repaymentType])

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
              <HomeIcon data-icon="inline-start" />
              Mortgage calculator
            </Badge>
            <h1 className="font-heading text-3xl font-semibold tracking-tight sm:text-4xl">
              Mortgage Repayments Calculator
            </h1>
            <p className="max-w-2xl text-muted-foreground">
              Estimate your monthly mortgage payments, total interest, and the
              full cost of your loan — before you commit.
            </p>
          </div>

          {/* Tabs */}
          <Tabs defaultValue="results">
            <TabsList>
              <TabsTrigger value="results">Results</TabsTrigger>
              <TabsTrigger value="amortisation">
                <FileSpreadsheet data-icon="inline-start" />
                Amortisation schedule
              </TabsTrigger>
            </TabsList>

            <TabsContent value="results">
              {/* Calculator grid */}
              <div className="grid grid-cols-1 gap-6 lg:grid-cols-5">
            {/* Input form */}
            <Card className="lg:col-span-3">
              <CardHeader>
                <CardTitle>Loan details</CardTitle>
                <CardDescription>
                  Enter your mortgage details to calculate repayments.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <FieldGroup>
                  <Field>
                    <FieldLabel htmlFor="home-price">Home price</FieldLabel>
                    <div className="relative">
                      <span className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-sm text-muted-foreground">
                        $
                      </span>
                      <Input
                        id="home-price"
                        type="number"
                        className="pl-6"
                        value={homePrice}
                        onChange={(e) => setHomePrice(e.target.value)}
                        min={0}
                        step={1000}
                      />
                    </div>
                  </Field>

                  <Field>
                    <FieldLabel htmlFor="deposit">Deposit</FieldLabel>
                    <div className="relative">
                      <span className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-sm text-muted-foreground">
                        $
                      </span>
                      <Input
                        id="deposit"
                        type="number"
                        className="pl-6"
                        value={deposit}
                        onChange={(e) => setDeposit(e.target.value)}
                        min={0}
                        step={1000}
                      />
                    </div>
                    <FieldDescription>
                      The amount you&apos;re putting down upfront.
                    </FieldDescription>
                  </Field>

                  <Field>
                    <FieldLabel htmlFor="interest-rate">
                      Interest rate
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
                    <FieldLabel>Loan term</FieldLabel>
                    <ToggleGroup
                      defaultValue={["25"]}
                      value={loanTerm}
                      onValueChange={(v) => {
                        if (v.length > 0) setLoanTerm(v)
                      }}
                      spacing={2}
                    >
                      {loanTermOptions.map((term) => (
                        <ToggleGroupItem key={term} value={term}>
                          {term} yr
                        </ToggleGroupItem>
                      ))}
                    </ToggleGroup>
                  </Field>

                  <Field>
                    <FieldLabel>Repayment type</FieldLabel>
                    <ToggleGroup
                      defaultValue={["repayment"]}
                      value={repaymentType}
                      onValueChange={(v) => {
                        if (v.length > 0) setRepaymentType(v)
                      }}
                      spacing={2}
                    >
                      <ToggleGroupItem value="repayment">
                        Repayment
                      </ToggleGroupItem>
                      <ToggleGroupItem value="interest-only">
                        Interest-only
                      </ToggleGroupItem>
                    </ToggleGroup>
                  </Field>
                </FieldGroup>
              </CardContent>
            </Card>

            {/* Results */}
            <div className="flex flex-col gap-6 lg:col-span-2">
              {/* Monthly payment */}
              <Card>
                <CardHeader>
                  <CardDescription>Monthly repayment</CardDescription>
                  <CardAction>
                    <CardTitle className="text-3xl font-semibold tabular-nums">
                      {currencyFormatterFull.format(results.monthlyPayment)}
                    </CardTitle>
                  </CardAction>
                </CardHeader>
              </Card>

              {/* Totals */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-base">Cost breakdown</CardTitle>
                </CardHeader>
                <CardContent className="flex flex-col gap-4">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">
                      Loan amount
                    </span>
                    <span className="font-medium tabular-nums">
                      {currencyFormatter.format(results.principal)}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">
                      Total interest
                    </span>
                    <span className="font-medium tabular-nums text-destructive">
                      {currencyFormatter.format(results.totalInterest)}
                    </span>
                  </div>
                  <div className="border-t pt-4">
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-medium">
                        Total amount paid
                      </span>
                      <span className="font-heading text-lg font-semibold tabular-nums">
                        {currencyFormatter.format(results.totalPaid)}
                      </span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Pie chart */}
              {results.chartData.length > 0 && results.totalInterest > 0 && (
                <Card>
                  <CardHeader>
                    <CardTitle className="text-base">
                      Principal vs interest
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
                                currencyFormatter.format(value as number)
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
                          content={<ChartLegendContent nameKey="name" />}
                        />
                      </PieChart>
                    </ChartContainer>
                  </CardContent>
                </Card>
              )}

                </div>
              </div>
            </TabsContent>

            <TabsContent value="amortisation">
              {repaymentType[0] === "repayment" &&
              results.schedule.length > 0 ? (
                <Card>
                  <CardHeader>
                    <CardTitle className="text-base">
                      Amortisation schedule
                    </CardTitle>
                    <CardDescription>
                      Month-by-month breakdown of payments for a{" "}
                      {loanTerm[0]}-year mortgage at{" "}
                      {interestRate || "0"}% interest.
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="flex flex-col gap-3">
                      <p className="text-sm text-muted-foreground">
                        Expand a year to see each monthly payment. The schedule
                        lists every month, which principal is repaid, and the
                        balance remaining — for a repayment mortgage the balance
                        reaches $0 at the end of the term.
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
                                    {group.rows.length} months
                                  </span>
                                </span>
                              </AccordionTrigger>
                              <AccordionContent>
                                <Table>
                                  <TableHeader>
                                    <TableRow>
                                      <TableHead className="w-20">
                                        Month
                                      </TableHead>
                                      <TableHead className="text-right">
                                        Payment
                                      </TableHead>
                                      <TableHead className="text-right">
                                        Interest
                                      </TableHead>
                                      <TableHead className="text-right">
                                        Principal
                                      </TableHead>
                                      <TableHead className="text-right">
                                        Balance
                                      </TableHead>
                                    </TableRow>
                                  </TableHeader>
                                  <TableBody>
                                    {group.rows.map((row) => (
                                      <TableRow key={row.month}>
                                        <TableCell className="font-medium tabular-nums">
                                          {row.month}
                                        </TableCell>
                                        <TableCell className="text-right tabular-nums">
                                          {currencyFormatterFull.format(
                                            row.payment
                                          )}
                                        </TableCell>
                                        <TableCell className="text-right tabular-nums">
                                          {currencyFormatterFull.format(
                                            row.interest
                                          )}
                                        </TableCell>
                                        <TableCell className="text-right tabular-nums">
                                          {currencyFormatterFull.format(
                                            row.principal
                                          )}
                                        </TableCell>
                                        <TableCell className="text-right tabular-nums">
                                          {currencyFormatterFull.format(
                                            row.balance
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
                          No amortisation schedule available
                        </EmptyTitle>
                        <EmptyDescription>
                          An amortisation schedule shows how each payment splits
                          between interest and principal as the loan balance
                          reduces. For interest-only mortgages the balance never
                          changes, so switch to a repayment mortgage to see the
                          month-by-month breakdown.
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
