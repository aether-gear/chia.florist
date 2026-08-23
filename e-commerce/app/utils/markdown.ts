// app/utils/markdown.ts

/**
 * A lightweight, safe markdown-to-HTML parser designed for structured instructions,
 * payment guides, and formatted backend markdown texts.
 */
export const renderMarkdown = (rawMarkdown: string): string => {
  if (!rawMarkdown) return ''

  // Normalize line breaks
  const normalized = rawMarkdown.replace(/\r\n/g, '\n').replace(/\r/g, '\n').trim()

  // Helper to format inline markdown spans (bold, italic, code, links)
  const formatInline = (text: string): string => {
    let out = text
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')

    // Inline code: `code`
    out = out.replace(/`([^`]+)`/g, '<code class="px-2 py-0.5 bg-emerald-50 text-emerald-900 font-mono text-xs font-bold rounded border border-emerald-200 select-all">$1</code>')

    // Bold: **text** or __text__
    out = out.replace(/\*\*([^*]+)\*\*/g, '<strong class="font-extrabold text-emerald-900">$1</strong>')
    out = out.replace(/__([^_]+)__/g, '<strong class="font-extrabold text-emerald-900">$1</strong>')

    // Italic: *text* or _text_
    out = out.replace(/\*([^*]+)\*/g, '<em class="italic text-gray-700">$1</em>')
    out = out.replace(/_([^_]+)_/g, '<em class="italic text-gray-700">$1</em>')

    // Links: [label](url)
    out = out.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer" class="text-emerald-700 hover:text-emerald-900 font-bold underline">$1 ↗</a>')

    return out
  }

  const lines = normalized.split('\n')
  const htmlChunks: string[] = []

  let inUl = false
  let inOl = false
  let paragraphBuffer: string[] = []

  const flushParagraph = () => {
    if (paragraphBuffer.length > 0) {
      const paragraphText = paragraphBuffer.map(l => formatInline(l)).join('<br/>')
      htmlChunks.push(`<p class="text-sm text-gray-700 leading-relaxed mb-3">${paragraphText}</p>`)
      paragraphBuffer = []
    }
  }

  const closeLists = () => {
    if (inUl) {
      htmlChunks.push('</ul>')
      inUl = false
    }
    if (inOl) {
      htmlChunks.push('</ol>')
      inOl = false
    }
  }

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    const trimmed = line.trim()

    // Blank line
    if (!trimmed) {
      flushParagraph()
      closeLists()
      continue
    }

    // Horizontal Rule: --- or *** or ___
    if (/^(\-{3,}|\*{3,}|_{3,})$/.test(trimmed)) {
      flushParagraph()
      closeLists()
      htmlChunks.push('<hr class="my-4 border-gray-200" />')
      continue
    }

    // Headings
    const h1Match = trimmed.match(/^#\s+(.+)$/)
    if (h1Match) {
      flushParagraph()
      closeLists()
      htmlChunks.push(`<h1 class="text-xl sm:text-2xl font-black text-gray-900 mt-4 mb-3 border-b border-emerald-100 pb-2">${formatInline(h1Match[1])}</h1>`)
      continue
    }

    const h2Match = trimmed.match(/^##\s+(.+)$/)
    if (h2Match) {
      flushParagraph()
      closeLists()
      htmlChunks.push(`<h2 class="text-lg sm:text-xl font-bold text-gray-900 mt-4 mb-2.5">${formatInline(h2Match[1])}</h2>`)
      continue
    }

    const h3Match = trimmed.match(/^###\s+(.+)$/)
    if (h3Match) {
      flushParagraph()
      closeLists()
      htmlChunks.push(`<h3 class="text-base sm:text-lg font-bold text-emerald-950 mt-3.5 mb-2">${formatInline(h3Match[1])}</h3>`)
      continue
    }

    const h4Match = trimmed.match(/^####\s+(.+)$/)
    if (h4Match) {
      flushParagraph()
      closeLists()
      htmlChunks.push(`<h4 class="text-sm sm:text-base font-bold text-gray-800 mt-3 mb-1.5">${formatInline(h4Match[1])}</h4>`)
      continue
    }

    const h5Match = trimmed.match(/^#####\s+(.+)$/)
    if (h5Match) {
      flushParagraph()
      closeLists()
      htmlChunks.push(`<h5 class="text-xs sm:text-sm font-bold text-gray-800 mt-2.5 mb-1">${formatInline(h5Match[1])}</h5>`)
      continue
    }

    // Blockquote: > text
    const bqMatch = trimmed.match(/^>\s*(.+)$/)
    if (bqMatch) {
      flushParagraph()
      closeLists()
      htmlChunks.push(`<blockquote class="border-l-4 border-emerald-500 bg-emerald-50/40 pl-3.5 py-1.5 my-2.5 rounded-r-xl text-xs text-gray-700 italic">${formatInline(bqMatch[1])}</blockquote>`)
      continue
    }

    // Unordered list item: - item or * item or + item
    const ulMatch = trimmed.match(/^[-*+]\s+(.+)$/)
    if (ulMatch) {
      flushParagraph()
      if (inOl) {
        htmlChunks.push('</ol>')
        inOl = false
      }
      if (!inUl) {
        htmlChunks.push('<ul class="list-disc pl-5 my-2.5 space-y-1 text-sm text-gray-700">')
        inUl = true
      }
      htmlChunks.push(`<li class="leading-relaxed font-medium">${formatInline(ulMatch[1])}</li>`)
      continue
    }

    // Ordered list item: 1. item, 2. item, etc.
    const olMatch = trimmed.match(/^(\d+)\.\s+(.+)$/)
    if (olMatch) {
      flushParagraph()
      if (inUl) {
        htmlChunks.push('</ul>')
        inUl = false
      }
      if (!inOl) {
        htmlChunks.push('<ol class="list-decimal pl-5 my-2.5 space-y-1.5 text-sm text-gray-700">')
        inOl = true
      }
      htmlChunks.push(`<li class="leading-relaxed font-medium">${formatInline(olMatch[2])}</li>`)
      continue
    }

    // Regular line in paragraph
    closeLists()
    paragraphBuffer.push(trimmed)
  }

  flushParagraph()
  closeLists()

  return htmlChunks.join('')
}
