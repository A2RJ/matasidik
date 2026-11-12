import { tv } from 'tailwind-variants';

export const buttonVariants = tv({
  base: [
    'inline-flex items-center justify-center gap-2 border-2 border-black px-4 font-mono text-xs font-bold uppercase tracking-[0.08em]',
    'transition-transform duration-100 active:translate-x-[2px] active:translate-y-[2px] active:shadow-none',
    'disabled:pointer-events-none disabled:opacity-50'
  ],
  variants: {
    variant: {
      default: 'bg-primary text-primary-foreground shadow-nb hover:-translate-y-0.5',
      secondary: 'bg-secondary text-secondary-foreground shadow-nb hover:-translate-y-0.5',
      outline: 'bg-card text-foreground shadow-nbsm hover:bg-primary',
      ghost: 'border-transparent bg-transparent text-foreground hover:border-black hover:bg-primary',
      destructive: 'bg-destructive text-destructive-foreground shadow-nb hover:-translate-y-0.5',
      success: 'bg-success text-success-foreground shadow-nb hover:-translate-y-0.5',
      dark: 'bg-foreground text-white shadow-nb hover:-translate-y-0.5'
    },
    size: {
      default: 'h-10',
      sm: 'h-8 px-3 text-[10px]',
      lg: 'h-12 px-6 text-sm',
      icon: 'h-10 w-10 px-0'
    }
  },
  defaultVariants: { variant: 'default', size: 'default' }
});
