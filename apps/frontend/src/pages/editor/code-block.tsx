import { useTheme } from "@/contexts/ThemeContext";
import { PrismAsyncLight as SyntaxHighlighter } from "react-syntax-highlighter";
import oneDark from "react-syntax-highlighter/dist/esm/styles/prism/one-dark";
import oneLight from "react-syntax-highlighter/dist/esm/styles/prism/one-light";

export default function CodeBlock({ code, language }: { code: string; language: string }) {
  const { theme } = useTheme();
  const lang = language && language.toLowerCase();

  console.log("CodeBlock", { code, language, lang, theme });
  return (
    <SyntaxHighlighter
      language={lang}
      style={theme === "dark" ? oneDark : oneLight}
    >
      {code}
    </SyntaxHighlighter>
  );
}
