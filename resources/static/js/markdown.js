function initializeMarkdownEditor() {
    let markdownEditorContainer = document.getElementById("markdownEditorContainer");
    let markdownEditor = document.getElementById("markdownEditor");

    if (markdownEditor) {
        markdownEditor.remove();
        markdownEditorContainer.innerHTML = "";
    }

    const newMarkdownEditor = document.createElement("textarea");
    newMarkdownEditor.id = "markdownEditor";
    markdownEditorContainer.appendChild(newMarkdownEditor);

    const simplemde = new SimpleMDE({
        element: newMarkdownEditor,
    });

    simplemde.codemirror.on("change", () => {
        markdownEditor = document.getElementById("inputMarkdown");
        markdownEditor.innerHTML = simplemde.value();
    });
}
