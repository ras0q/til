export function jsx(
  tag: string,
  { children, ...props }: Record<string, string>,
): string {
  const attributes = Object.entries(props)
    .map(([key, value]) => ` ${key}="${value}"`)
    .join("");
  const innerHTML = Array.isArray(children) ? children.join("") : children;
  return `<${tag}${attributes}>${innerHTML}</${tag}>`;
}

export { jsx as jsxs };
