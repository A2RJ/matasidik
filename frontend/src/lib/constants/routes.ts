export const ROUTES = {
  dashboard: '/',
  upload: '/upload',
  submissions: '/submissions',
  submission: (id: string) => `/submissions/${id}`,
  master: '/master',
  factors: '/factors',
  audit: '/audit'
} as const;
