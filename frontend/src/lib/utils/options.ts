export function toOptions<T extends Record<string, any>>(
  items: T[],
  valueKey: keyof T = 'id' as keyof T,
  labelKey: keyof T = 'nama' as keyof T
): { value: string; label: string }[] {
  return items.map((item) => ({ value: String(item[valueKey]), label: String(item[labelKey]) }));
}
