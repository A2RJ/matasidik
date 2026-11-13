export const ROUTES = {
  dashboard: '/',
  userGuide: '/user-guide',
  upload: '/upload',
  submissions: '/submissions',
  submission: (id: string) => `/submissions/${id}`,
  master: '/master',
  factors: '/factors',
  audit: '/audit'
} as const;
