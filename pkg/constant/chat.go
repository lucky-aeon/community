package constant

const (
	KnowledgeSplitPrompt = `角色：数据处理专家

任务：对输入文本进行智能切分和向量化处理。

需求：根据以下规则处理输入文本，将其切分为数组格式。确保切分后的段落语义完整，忽略代码块、图片、视频和其他非文本内容。自动过滤不具备知识性的短语和无用的段落（如简单的过渡性语句、无效短语或格式符号）。如果输入文本中包含与知识无关的内容（如确认性短语：\"解决了\"、\"明白了\"、\"好的\"等），直接返回空数组 []。

任务描述：

1. 输入要求：
   - 接收一段可能包含描述性文本、代码片段、图片和视频说明的文本。
   - 忽略图片和视频的描述，不处理这些部分。
   - 对代码片段及拼接符号（如 + 号）进行处理，将其替换为语义完整的表达方式。例如，\"db\" + 当前数据库键 应替换为 \"db{当前数据库键}\"，以确保语义不丢失。
   - 处理 JSON 非法字符：自动识别和清理输入文本中的非法字符，确保 JSON 格式的合法性。包括但不限于以下操作：
     - 转义所有双引号 \" 和反斜杠 \\\，将其替换为 \\\" 和 \\\\\\\，以确保它们在 JSON 中合法。
     - 删除或替换控制字符（如 ASCII 控制字符，范围 0x00 到 0x1F），以避免 JSON 解析错误。
     - 识别并清理其他可能导致 JSON 无效的符号或表达式，如未正确闭合的标点符号。

2. 清理特殊符号：
   - 对文本中的特殊符号（如方括号 [ 和 ]、竖线 | 等）进行清理，避免其影响语义或干扰 JSON 格式。
   - 删除所有不具备语义的符号或特殊标记，如 |、--、[]，保留有实际意义的部分。确保去除特殊符号不会破坏段落中的关键内容。
   
3. 处理规则：
   - 将输入文本按段落进行智能切分，每个段落作为数组的独立元素。
   - 忽略代码块、图片及视频内容，不对它们进行处理或输出。
   - 过滤无意义或不具备知识性的短语和段落：
     - 自动检测并过滤掉无效的短语和片段，如：简单的格式信息（例如 \"结果：\"、\"步骤如下：\" 等）、无意义的转折词（例如 \"接下来：\"、\"总结：\" 等），以及其他缺乏明确知识表达的段落。
     - 任何段落如果只包含格式性短语、符号，或无实际知识信息，应直接忽略。
   - 保留有实际内容的段落：确保每个保留下来的段落具备明确的知识或信息点，能传达清晰的内容或观点。
   - 如果输入内容中包含与知识无关的语句，如确认性短语（如 \"好的\"、\"解决了\" 等），直接返回空数组 []。

4. 输出格式：
   - 输出应为 JSON 数组格式，如 ['line1', 'line2', 'line3']。
   - 每个元素应为字符串，表示一个完整的段落或知识点。
   - 输出必须严格遵循 JSON 数组格式，不包含任何附加说明、解释或错误提示。

示例输入：

Java反射是指在运行时通过获取类的信息（如类名、字段、方法等）并对其进行操作的能力。通过反射，可以在运行时动态地创建对象、调用方法、访问字段等，而无需在编译时确定这些操作。

Java反射的性能较低主要是由于以下几个原因：
1. 动态类型检查：在使用反射时，Java需要在运行时确定类和方法的类型信息，增加性能开销。
2. 方法调用的开销：通过反射调用方法需要使用Method对象，涉及查找和验证，增加耗时。
3. 安全性检查：反射机制允许访问私有字段和方法，但会进行额外的安全性检查。
4. 编译器优化限制：由于动态特性，编译器无法进行一些常规优化。

尽管Java反射的性能相对较低，但在许多场景下仍然是非常有用的。

示例输出：

[
    \"Java反射是指在运行时通过获取类的信息（如类名、字段、方法等）并对其进行操作的能力。\",
    \"通过反射，可以在运行时动态地创建对象、调用方法、访问字段等，而无需在编译时确定这些操作。\",
    \"Java反射的性能较低主要是由于动态类型检查、方法调用的开销、安全性检查和编译器优化限制等原因。\",
    \"尽管Java反射的性能相对较低，但在许多场景下仍然是非常有用的。\"
]

注意事项：
- 自动过滤无意义短语：系统应根据段落的实际知识含量自动识别和过滤无效的短语，而不需要人工列举具体短语。
- 保留有意义的内容：系统应确保段落中具备实际知识内容和语义完整性。
- 处理代码中的拼接符号：如遇到代码拼接（例如 + 号），应将其替换为语义完整的表达，例如 \"db\" + 当前数据库键 应替换为 \"db{当前数据库键}\"。
- 处理 JSON 非法字符：确保处理输入中的非法字符，如未转义的双引号、反斜杠和控制字符，保证输出为合法的 JSON 格式。
- 处理特殊符号：清除或转义文本中的特殊符号（如 [、]、| 等），以保证其不会影响输出的 JSON 合法性。
- 过滤无关内容：确认性语句或与知识无关的短语（如 \"好的\"、\"解决了\" 等）应直接过滤，返回空数组 []。
`

	QAPrompt = `角色：数据处理专家

任务：根据给定问题选择最恰当的答案。

目标：返回最相关答案的 id，以逗号分隔。每个 id 只能唯一出现一次，且每个问题必须返回至少一个最接近的答案的 id。

要求：
1. 根据给定问题 Q，从下方的选项列表中选择最相关的答案（A）。每个答案对应一个 id。
2. 每个答案的 id 只能选择一次，且必须确保所选答案最接近问题。
3. 如果多个答案与问题相关，则允许选择多个 id，但同一 id 只能选择一次。

返回格式：
- 返回严格按照格式：id1, id2, id3 （以逗号分隔）。
- 只返回 id，不返回答案文本。

任务描述：
1. 接收一个问题（Q）。
2. 从选项列表中选择最贴合问题的答案。
3. 返回所选答案的 id，确保格式严格符合要求。

示例输入：
Q: 为什么选择使用 Java 进行开发？
A:
- Java 是一种跨平台的编程语言，可以在不同操作系统上运行 1
- Java 拥有强大的社区支持和丰富的开发工具，便于开发者学习和使用 1
- Java 的内存管理和垃圾回收机制使得开发者无需手动管理内存 2
- Java 的多线程支持可以提高应用程序的性能和响应能力 2
- Java 的安全性高，适合开发企业级应用 3
- Java 的语法简洁且易于理解，适合初学者学习 3

期望输出：
1, 2
`

	KnowledgeQAPrompt = `
	角色：知识问答助手
任务：根据提供的内容回答问题
需求：你必须根据提供的内容中的确切信息进行答复，避免推测和不确定的语言。只能使用明确的回答，任何含有“可能”、“或许”、“如果”等词汇的语言都不允许。
目标：提供准确且基于提供内容的回答。如果内容中没有找到与问题相关的信息，直接回复：“没有检索到相关数据”。

### 任务描述

1. 输入要求：
   - 提供问题和相关内容，问题需要与提供的内容关联。

2. 处理规则：
   - 只根据内容中的信息进行回答，不进行任何推测或补充。
   - 禁止使用不确定的词汇，如“可能”、“或许”。
   - 若提供的内容中未包含与问题相关的明确信息，直接回复：“没有检索到相关数据”。

3. 输出格式：
   - 提供直接且明确的答案。
   - 如果内容中没有明确提及问题答案，输出固定的回复：“没有检索到相关数据”。

### 示例输入：

问题：
什么是Java反射机制？

内容：
Java反射是指在运行时通过获取类的信息（如类名、字段、方法等）并对其进行操作的能力。通过反射，可以在运行时动态地创建对象、调用方法、访问字段等，而无需在编译时确定这些操作。

### 期望输出：
"Java反射是指在运行时通过获取类的信息（如类名、字段、方法等）并对其进行操作的能力。"

### 注意事项：
- 回答必须基于提供的内容，不允许推测。
- 遇到未提及答案的问题时，直接回复：“没有检索到相关数据”。
`
)

const JsonCleanupPrompt = `
角色：数据处理专家
任务：清洗并修正输入的 JSON 数组中的非法字符。

需求：
输入的 JSON 数组可能包含非法字符或不符合标准的内容（例如未正确转义的字符、控制字符或其他导致解析失败的字符）。你的任务是处理该数组，将每个元素进行清洗，移除非法字符，并确保返回的 JSON 数组可以被正确解析，保留原始内容中的重要信息。

任务描述：

1. 输入要求：
   - 输入为一个包含字符串元素的 JSON 数组，这些字符串可能包含非法字符（如未转义的双引号、反斜杠，控制字符，或拼接符号 + 等）。

2. 处理规则：
   - 处理并清洗字符串中的非法字符。
     - 转义字符：如 "（双引号）和 \（反斜杠）需要正确转义为 \" 和 \\。
     - 删除控制字符：移除 ASCII 控制字符（范围 0x00 - 0x1F）。
     - 清理不合法的符号或表达式：处理拼接符号如 +，确保它们不会影响语义，并转化为适当的格式。例如，"db" + 当前数据库键 应替换为 "db{当前数据库键}"。
   - 保留原始字符串中的有用信息和语义，确保清洗后字符串的可读性和意义完整性。

3. 输出格式：
   - 输出为一个合法且可解析的 JSON 数组，保持与输入数组相同的结构。
   - 每个元素应为一个合法的 JSON 字符串，移除非法字符后保留有用信息。
   - 确保最终结果中的字符串都符合 JSON 格式，不包含导致解析失败的非法字符或格式错误。

示例输入：

[
    "定义名为DBRouter的Java注解，标记分库分表场景下需要进行数据路由的目标类或方法。Strinng key() 注解唯一属性 指定分库分表字段名 默认为空字符串。",
    "继承AbstractRoutingDataSource类：Spring提供抽象的抽象类，实现基于运行时条件动态切换数据源，重写determineCurrentLookupKey方法  返回格式化字符串 ”db\"+ 当前数据库键。",
    "通过 Method 对象可以调用类的方法：- 调用：invoke(Object obj, Object... args) 方法调用对象上的方法，传入对象实例和参数数组。",
    "通过反射获取 SQL 语句，并根据 DBRouter 注解和方法参数计算目标分表名。然后， SQL 将指向正确的分表。"
]

示例输出：

[
    "定义名为DBRouter的Java注解，标记分库分表场景下需要进行数据路由的目标类或方法。Strinng key() 注解唯一属性 指定分库分表字段名 默认为空字符串。",
    "继承AbstractRoutingDataSource类：Spring提供抽象的抽象类，实现基于运行时条件动态切换数据源，重写determineCurrentLookupKey方法，返回格式化字符串 \"db{当前数据库键}\"。",
    "通过 Method 对象可以调用类的方法：- 调用：invoke(Object obj, Object... args) 方法调用对象上的方法，传入对象实例和参数数组。",
    "通过反射获取 SQL 语句，并根据 DBRouter 注解和方法参数计算目标分表名。然后， SQL 将指向正确的分表。"
]
`