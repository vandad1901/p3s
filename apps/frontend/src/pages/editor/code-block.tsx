import { useTheme } from "@/contexts/ThemeContext";
import { PrismAsyncLight as SyntaxHighlighter } from "react-syntax-highlighter";
import oneDark from "react-syntax-highlighter/dist/esm/styles/prism/one-dark";
import oneLight from "react-syntax-highlighter/dist/esm/styles/prism/one-light";

export default function CodeBlock({ code, language }: { code: string; language: string }) {
  const { theme } = useTheme();
  const lang = language && language.toLowerCase();

  const style = theme === "dark" ? oneDark : oneLight;

  return (
    <SyntaxHighlighter
      language={lang}
      style={style}
      customStyle={{ borderRadius: "0 0.3rem 0.3rem 0.3rem" }}
      showLineNumbers={true}
    >
      {code}
    </SyntaxHighlighter>
  );
}
