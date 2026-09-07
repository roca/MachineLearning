export function Footer() {
  return (
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
  )
}
