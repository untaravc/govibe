<template>
  <div class="rounded-2xl border border-slate-200 bg-white">
    <div class="flex flex-wrap items-center gap-1 border-b border-slate-200 px-3 py-2">
      <button
        type="button"
        class="rounded-lg px-2 py-1 text-sm font-medium text-slate-700 hover:bg-slate-100"
        :class="editor?.isActive('bold') ? 'bg-slate-100 text-slate-900' : ''"
        @click="editor?.chain().focus().toggleBold().run()"
      >
        Bold
      </button>
      <button
        type="button"
        class="rounded-lg px-2 py-1 text-sm font-medium text-slate-700 hover:bg-slate-100"
        :class="editor?.isActive('italic') ? 'bg-slate-100 text-slate-900' : ''"
        @click="editor?.chain().focus().toggleItalic().run()"
      >
        Italic
      </button>
      <button
        type="button"
        class="rounded-lg px-2 py-1 text-sm font-medium text-slate-700 hover:bg-slate-100"
        :class="editor?.isActive('bulletList') ? 'bg-slate-100 text-slate-900' : ''"
        @click="editor?.chain().focus().toggleBulletList().run()"
      >
        Bullet
      </button>
      <button
        type="button"
        class="rounded-lg px-2 py-1 text-sm font-medium text-slate-700 hover:bg-slate-100"
        @click="editor?.chain().focus().setParagraph().run()"
      >
        Paragraph
      </button>
    </div>

    <div class="prose prose-slate max-w-none px-4 py-3">
      <EditorContent :editor="editor" />
    </div>
  </div>
</template>

<script setup>
import { onBeforeUnmount, watch } from "vue";
import { EditorContent, useEditor } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";

const props = defineProps({
  modelValue: {
    type: String,
    default: ""
  }
});

const emit = defineEmits(["update:modelValue"]);

const editor = useEditor({
  content: props.modelValue || "",
  extensions: [StarterKit],
  editorProps: {
    attributes: {
      class: "min-h-[160px] outline-none"
    }
  },
  onUpdate: ({ editor }) => {
    emit("update:modelValue", editor.getHTML());
  }
});

watch(
  () => props.modelValue,
  (next) => {
    if (!editor.value) return;
    const current = editor.value.getHTML();
    const desired = next || "";
    if (current !== desired) editor.value.commands.setContent(desired, false);
  }
);

onBeforeUnmount(() => {
  editor.value?.destroy();
});
</script>

