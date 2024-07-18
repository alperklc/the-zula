import React from "react";
import { BlockTypeSelect, BoldItalicUnderlineToggles, CodeToggle, headingsPlugin, InsertTable, InsertThematicBreak, listsPlugin, ListsToggle, MDXEditor, MDXEditorMethods, tablePlugin, thematicBreakPlugin, toolbarPlugin } from "@mdxeditor/editor";
import '@mdxeditor/editor/style.css'

import styles from "./index.module.css";

const MdEditor = (props: { className?: string; content: string; onChange: (_: string) => void }) => {
    const editorRef = React.useRef<MDXEditorMethods>()

    React.useEffect(() => {
        if (props.content?.length) {
            editorRef.current?.setMarkdown(props.content)
        }
    }, [props.content])

    return (<MDXEditor
        contentEditableClassName={`${props.className} ${styles.container}`}
        markdown={props.content}
        onChange={props.onChange}
        ref={editorRef as React.LegacyRef<MDXEditorMethods>}
        plugins={[
            headingsPlugin(),
            listsPlugin(),
            tablePlugin(),
            thematicBreakPlugin(),
            toolbarPlugin({
                toolbarContents: () => (
                    <>
                        <BoldItalicUnderlineToggles />
                        <CodeToggle />
                        <ListsToggle />
                        <BlockTypeSelect />
                        <InsertTable />
                        <InsertThematicBreak />
                    </>
                )
            })

        ]}
    />
    );
}
export default MdEditor;
