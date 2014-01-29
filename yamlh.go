package candiedyaml

import (
	"io"
)

/** The version directive data. */
type yaml_version_directive_t struct {
	major int // The major version number
	minor int // The minor version number
}

/** The tag directive data. */
type yaml_tag_directive_t struct {
	handle []byte // The tag handle
	prefix []byte // The tag prefix
}

/** The stream encoding. */
type yaml_encoding_t int

const (
	/** Let the parser choose the encoding. */
	YAML_ANY_ENCODING yaml_encoding_t = iota
	/** The defau lt UTF-8 encoding. */
	YAML_UTF8_ENCODING
	/** The UTF-16-LE encoding with BOM. */
	YAML_UTF16LE_ENCODING
	/** The UTF-16-BE encoding with BOM. */
	YAML_UTF16BE_ENCODING
)

/** Line break types. */
type yaml_break_t int

const (
	YAML_ANY_BREAK  yaml_break_t = iota /** Let the parser choose the break type. */
	YAML_CR_BREAK                       /** Use CR for line breaks (Mac style). */
	YAML_LN_BREAK                       /** Use LN for line breaks (Unix style). */
	YAML_CRLN_BREAK                     /** Use CR LN for line breaks (DOS style). */
)

/** Many bad things could happen with the parser and emitter. */
type yaml_error_type_t int

const (
	/** No error is produced. */
	YAML_NO_ERROR yaml_error_type_t = iota

	/** Cannot allocate or reallocate a block of memory. */
	YAML_MEMORY_ERROR

	/** Cannot read or decode the input stream. */
	YAML_READER_ERROR
	/** Cannot scan the input stream. */
	YAML_SCANNER_ERROR
	/** Cannot parse the input stream. */
	YAML_PARSER_ERROR
	/** Cannot compose a YAML document. */
	YAML_COMPOSER_ERROR

	/** Cannot write to the output stream. */
	YAML_WRITER_ERROR
	/** Cannot emit a YAML stream. */
	YAML_EMITTER_ERROR
)

/** The pointer position. */
type yaml_mark_t struct {
	/** The position index. */
	index int

	/** The position line. */
	line int

	/** The position column. */
	column int
}

/** @} */

/**
 * @defgroup styles Node Styles
 * @{
 */

type yaml_style_t int

/** Scalar styles. */
type yaml_scalar_style_t yaml_style_t

const (
	/** Let the emitter choose the style. */
	YAML_ANY_SCALAR_STYLE yaml_scalar_style_t = iota

	/** The plain scalar style. */
	YAML_PLAIN_SCALAR_STYLE

	/** The single-quoted scalar style. */
	YAML_SINGLE_QUOTED_SCALAR_STYLE
	/** The double-quoted scalar style. */
	YAML_DOUBLE_QUOTED_SCALAR_STYLE

	/** The literal scalar style. */
	YAML_LITERAL_SCALAR_STYLE
	/** The folded scalar style. */
	YAML_FOLDED_SCALAR_STYLE
)

/** Sequence styles. */
type yaml_sequence_style_t yaml_style_t

const (
	/** Let the emitter choose the style. */
	YAML_ANY_SEQUENCE_STYLE yaml_sequence_style_t = iota

	/** The block sequence style. */
	YAML_BLOCK_SEQUENCE_STYLE
	/** The flow sequence style. */
	YAML_FLOW_SEQUENCE_STYLE
)

/** Mapping styles. */
type yaml_mapping_style_t yaml_style_t

const (
	/** Let the emitter choose the style. */
	YAML_ANY_MAPPING_STYLE yaml_mapping_style_t = iota

	/** The block mapping style. */
	YAML_BLOCK_MAPPING_STYLE
	/** The flow mapping style. */
	YAML_FLOW_MAPPING_STYLE

/*    YAML_FLOW_SET_MAPPING_STYLE   */
)

/** @} */

/**
 * @defgroup tokens Tokens
 * @{
 */

/** Token types. */
type yaml_token_type_t int

const (
	/** An empty token. */
	YAML_NO_TOKEN yaml_token_type_t = iota

	/** A STREAM-START token. */
	YAML_STREAM_START_TOKEN
	/** A STREAM-END token. */
	YAML_STREAM_END_TOKEN

	/** A VERSION-DIRECTIVE token. */
	YAML_VERSION_DIRECTIVE_TOKEN
	/** A TAG-DIRECTIVE token. */
	YAML_TAG_DIRECTIVE_TOKEN
	/** A DOCUMENT-START token. */
	YAML_DOCUMENT_START_TOKEN
	/** A DOCUMENT-END token. */
	YAML_DOCUMENT_END_TOKEN

	/** A BLOCK-SEQUENCE-START token. */
	YAML_BLOCK_SEQUENCE_START_TOKEN
	/** A BLOCK-SEQUENCE-END token. */
	YAML_BLOCK_MAPPING_START_TOKEN
	/** A BLOCK-END token. */
	YAML_BLOCK_END_TOKEN

	/** A FLOW-SEQUENCE-START token. */
	YAML_FLOW_SEQUENCE_START_TOKEN
	/** A FLOW-SEQUENCE-END token. */
	YAML_FLOW_SEQUENCE_END_TOKEN
	/** A FLOW-MAPPING-START token. */
	YAML_FLOW_MAPPING_START_TOKEN
	/** A FLOW-MAPPING-END token. */
	YAML_FLOW_MAPPING_END_TOKEN

	/** A BLOCK-ENTRY token. */
	YAML_BLOCK_ENTRY_TOKEN
	/** A FLOW-ENTRY token. */
	YAML_FLOW_ENTRY_TOKEN
	/** A KEY token. */
	YAML_KEY_TOKEN
	/** A VALUE token. */
	YAML_VALUE_TOKEN

	/** An ALIAS token. */
	YAML_ALIAS_TOKEN
	/** An ANCHOR token. */
	YAML_ANCHOR_TOKEN
	/** A TAG token. */
	YAML_TAG_TOKEN
	/** A SCALAR token. */
	YAML_SCALAR_TOKEN
)

/** The token structure. */
type yaml_token_t struct {

	/** The token type. */
	token_type yaml_token_type_t

	/** The token data. */
	/** The stream start (for @c YAML_STREAM_START_TOKEN). */
	encoding yaml_encoding_t

	/** The alias (for @c YAML_ALIAS_TOKEN, YAML_ANCHOR_TOKEN, YAML_SCALAR_TOKEN,YAML_TAG_TOKEN ). */
	/** The anchor (for @c ). */
	/** The scalar value (for @c ). */
	value []byte

	/** The tag suffix. */
	suffix []byte

	/** The scalar value (for @c YAML_SCALAR_TOKEN). */
	/** The scalar style. */
	style yaml_scalar_style_t

	/** The version directive (for @c YAML_VERSION_DIRECTIVE_TOKEN). */
	version_directive yaml_version_directive_t

	/** The tag directive (for @c YAML_TAG_DIRECTIVE_TOKEN). */
	prefix []byte

	/** The beginning of the token. */
	start_mark yaml_mark_t
	/** The end of the token. */
	end_mark yaml_mark_t
	
	major, minor int
}

/**
 * @defgroup events Events
 * @{
 */

/** Event types. */
type yaml_event_type_t int

const (
	/** An empty event. */
	YAML_NO_EVENT yaml_event_type_t = iota

	/** A STREAM-START event. */
	YAML_STREAM_START_EVENT
	/** A STREAM-END event. */
	YAML_STREAM_END_EVENT

	/** A DOCUMENT-START event. */
	YAML_DOCUMENT_START_EVENT
	/** A DOCUMENT-END event. */
	YAML_DOCUMENT_END_EVENT

	/** An ALIAS event. */
	YAML_ALIAS_EVENT
	/** A SCALAR event. */
	YAML_SCALAR_EVENT

	/** A SEQUENCE-START event. */
	YAML_SEQUENCE_START_EVENT
	/** A SEQUENCE-END event. */
	YAML_SEQUENCE_END_EVENT

	/** A MAPPING-START event. */
	YAML_MAPPING_START_EVENT
	/** A MAPPING-END event. */
	YAML_MAPPING_END_EVENT
)

/** The event structure. */
type yaml_event_t struct {

	/** The event type. */
	event_type yaml_event_type_t

	/** The stream parameters (for @c YAML_STREAM_START_EVENT). */
	encoding yaml_encoding_t

	/** The document parameters (for @c YAML_DOCUMENT_START_EVENT). */
	version_directive *yaml_version_directive_t

	/** The beginning and end of the tag directives list. */
	tag_directives []yaml_tag_directive_t

	/** The document parameters (for @c YAML_DOCUMENT_START_EVENT, YAML_DOCUMENT_END_EVENT, YAML_SEQUENCE_START_EVENT,YAML_MAPPING_START_EVENT). */
	/** Is the document indicator implicit? */
	implicit bool

	/** The alias parameters (for @c YAML_ALIAS_EVENT,YAML_SCALAR_EVENT, YAML_SEQUENCE_START_EVENT, YAML_MAPPING_START_EVENT). */
	/** The anchor. */
	anchor []byte

	/** The scalar parameters (for @c YAML_SCALAR_EVENT,YAML_SEQUENCE_START_EVENT, YAML_MAPPING_START_EVENT). */
	/** The tag. */
	tag []byte
	/** The scalar value. */
	value []byte

	/** Is the tag optional for the plain style? */
	plain_implicit bool
	/** Is the tag optional for any non-plain style? */
	quoted_implicit bool

	/** The sequence parameters (for @c YAML_SEQUENCE_START_EVENT, YAML_MAPPING_START_EVENT). */
	/** The sequence style. */
	/** The scalar style. */
	style yaml_style_t

	/** The beginning of the event. */
	start_mark, end_mark yaml_mark_t
}

/**
 * @defgroup nodes Nodes
 * @{
 */

const (
	/** The tag @c !!null with the only possible value: @c null. */
	YAML_NULL_TAG = "tag:yaml.org,2002:null"
	/** The tag @c !!bool with the values: @c true and @c falce. */
	YAML_BOOL_TAG = "tag:yaml.org,2002:bool"
	/** The tag @c !!str for string values. */
	YAML_STR_TAG = "tag:yaml.org,2002:str"
	/** The tag @c !!int for integer values. */
	YAML_INT_TAG = "tag:yaml.org,2002:int"
	/** The tag @c !!float for float values. */
	YAML_FLOAT_TAG = "tag:yaml.org,2002:float"
	/** The tag @c !!timestamp for date and time values. */
	YAML_TIMESTAMP_TAG = "tag:yaml.org,2002:timestamp"

	/** The tag @c !!seq is used to denote sequences. */
	YAML_SEQ_TAG = "tag:yaml.org,2002:seq"
	/** The tag @c !!map is used to denote mapping. */
	YAML_MAP_TAG = "tag:yaml.org,2002:map"

	/** The default scalar tag is @c !!str. */
	YAML_DEFAULT_SCALAR_TAG = YAML_STR_TAG
	/** The default sequence tag is @c !!seq. */
	YAML_DEFAULT_SEQUENCE_TAG = YAML_SEQ_TAG
	/** The default mapping tag is @c !!map. */
	YAML_DEFAULT_MAPPING_TAG = YAML_MAP_TAG
)

/** Node types. */
type yaml_node_type_t int

const (
	/** An empty node. */
	YAML_NO_NODE yaml_node_type_t = iota

	/** A scalar node. */
	YAML_SCALAR_NODE
	/** A sequence node. */
	YAML_SEQUENCE_NODE
	/** A mapping node. */
	YAML_MAPPING_NODE
)

/** An element of a sequence node. */
type yaml_node_item_t int

/** An element of a mapping node. */
type yaml_node_pair_t struct {
	/** The key of the element. */
	key int
	/** The value of the element. */
	value int
}

/** The node structure. */
type yaml_node_t struct {

	/** The node type. */
	node_type yaml_node_type_t

	/** The node tag. */
	tag []byte

	/** The scalar parameters (for @c YAML_SCALAR_NODE). */
	scalar struct {
		/** The scalar value. */
		value []byte
		/** The scalar style. */
		style yaml_scalar_style_t
	}

	/** The sequence parameters (for @c YAML_SEQUENCE_NODE). */
	sequence struct {
		/** The stack of sequence items. */
		items []yaml_node_item_t
		/** The sequence style. */
		style yaml_sequence_style_t
	}

	/** The mapping parameters (for @c YAML_MAPPING_NODE). */
	mapping struct {
		/** The stack of mapping pairs (key, value). */
		pairs []yaml_node_pair_t
		/** The mapping style. */
		style yaml_mapping_style_t
	}

	/** The beginning of the node. */
	start_mark yaml_mark_t
	/** The end of the node. */
	end_mark yaml_mark_t
}

/** The document structure. */
type yaml_document_t struct {

	/** The document nodes. */
	nodes []yaml_node_t

	/** The version directive. */
	version_directive *yaml_version_directive_t

	/** The list of tag directives. */
	tags []yaml_tag_directive_t

	/** Is the document start indicator implicit? */
	start_implicit bool
	/** Is the document end indicator implicit? */
	end_implicit bool

	/** The beginning of the document. */
	start_mark yaml_mark_t
	/** The end of the document. */
	end_mark yaml_mark_t
}

/**
 * The prototype of a read handler.
 *
 * The read handler is called when the parser needs to read more bytes from the
 * source.  The handler should write not more than @a size bytes to the @a
 * buffer.  The number of written bytes should be set to the @a length variable.
 *
 * @param[in,out]   data        A pointer to an application data specified by
 *                              yaml_parser_set_input().
 * @param[out]      buffer      The buffer to write the data from the source.
 * @param[in]       size        The size of the buffer.
 * @param[out]      size_read   The actual number of bytes read from the source.
 *
 * @returns On success, the handler should return @c 1.  If the handler failed,
 * the returned value should be @c 0.  On EOF, the handler should set the
 * @a size_read to @c 0 and return @c 1.
 */

type yaml_read_handler_t func(parser *yaml_parser_t, buffer []byte) (n int, err error)

/**
 * This structure holds information about a potential simple key.
 */

type yaml_simple_key_t struct {
	/** Is a simple key possible? */
	possible bool

	/** Is a simple key required? */
	required bool

	/** The number of the token. */
	token_number int

	/** The position mark. */
	mark yaml_mark_t
}

/**
 * The states of the parser.
 */
type yaml_parser_state_t int

const (
	/** Expect STREAM-START. */
	YAML_PARSE_STREAM_START_STATE yaml_parser_state_t = iota
	/** Expect the beginning of an implicit document. */
	YAML_PARSE_IMPLICIT_DOCUMENT_START_STATE
	/** Expect DOCUMENT-START. */
	YAML_PARSE_DOCUMENT_START_STATE
	/** Expect the content of a document. */
	YAML_PARSE_DOCUMENT_CONTENT_STATE
	/** Expect DOCUMENT-END. */
	YAML_PARSE_DOCUMENT_END_STATE
	/** Expect a block node. */
	YAML_PARSE_BLOCK_NODE_STATE
	/** Expect a block node or indentless sequence. */
	YAML_PARSE_BLOCK_NODE_OR_INDENTLESS_SEQUENCE_STATE
	/** Expect a flow node. */
	YAML_PARSE_FLOW_NODE_STATE
	/** Expect the first entry of a block sequence. */
	YAML_PARSE_BLOCK_SEQUENCE_FIRST_ENTRY_STATE
	/** Expect an entry of a block sequence. */
	YAML_PARSE_BLOCK_SEQUENCE_ENTRY_STATE
	/** Expect an entry of an indentless sequence. */
	YAML_PARSE_INDENTLESS_SEQUENCE_ENTRY_STATE
	/** Expect the first key of a block mapping. */
	YAML_PARSE_BLOCK_MAPPING_FIRST_KEY_STATE
	/** Expect a block mapping key. */
	YAML_PARSE_BLOCK_MAPPING_KEY_STATE
	/** Expect a block mapping value. */
	YAML_PARSE_BLOCK_MAPPING_VALUE_STATE
	/** Expect the first entry of a flow sequence. */
	YAML_PARSE_FLOW_SEQUENCE_FIRST_ENTRY_STATE
	/** Expect an entry of a flow sequence. */
	YAML_PARSE_FLOW_SEQUENCE_ENTRY_STATE
	/** Expect a key of an ordered mapping. */
	YAML_PARSE_FLOW_SEQUENCE_ENTRY_MAPPING_KEY_STATE
	/** Expect a value of an ordered mapping. */
	YAML_PARSE_FLOW_SEQUENCE_ENTRY_MAPPING_VALUE_STATE
	/** Expect the and of an ordered mapping entry. */
	YAML_PARSE_FLOW_SEQUENCE_ENTRY_MAPPING_END_STATE
	/** Expect the first key of a flow mapping. */
	YAML_PARSE_FLOW_MAPPING_FIRST_KEY_STATE
	/** Expect a key of a flow mapping. */
	YAML_PARSE_FLOW_MAPPING_KEY_STATE
	/** Expect a value of a flow mapping. */
	YAML_PARSE_FLOW_MAPPING_VALUE_STATE
	/** Expect an empty value of a flow mapping. */
	YAML_PARSE_FLOW_MAPPING_EMPTY_VALUE_STATE
	/** Expect nothing. */
	YAML_PARSE_END_STATE
)

/**
 * This structure holds aliases data.
 */

type yaml_alias_data_t struct {
	/** The anchor. */
	anchor []byte
	/** The node id. */
	index int
	/** The anchor mark. */
	mark yaml_mark_t
}

/**
 * The parser structure.
 *
 * All members are internal.  Manage the structure using the @c yaml_parser_
 * family of functions.
 */

type yaml_parser_t struct {

	/**
	 * @name Error handling
	 * @{
	 */

	/** Error type. */
	error yaml_error_type_t
	/** Error description. */
	problem string
	/** The byte about which the problem occured. */
	problem_offset int
	/** The problematic value (@c -1 is none). */
	problem_value int
	/** The problem position. */
	problem_mark yaml_mark_t
	/** The error context. */
	context string
	/** The context position. */
	context_mark yaml_mark_t

	/**
	 * @}
	 */

	/**
	 * @name Reader stuff
	 * @{
	 */

	/** Read handler. */
	read_handler yaml_read_handler_t

	/** File input data. */
	input_file io.Reader
	input      []byte
	input_pos  int

	/** EOF flag */
	eof bool

	/** The working buffer. */
	buffer     []byte
	buffer_pos int

	/* The number of unread characters in the buffer. */
	unread int

	/** The raw buffer. */
	raw_buffer     []byte
	raw_buffer_pos int

	/** The input encoding. */
	encoding yaml_encoding_t

	/** The offset of the current position (in bytes). */
	offset int

	/** The mark of the current position. */
	mark yaml_mark_t

	/**
	 * @}
	 */

	/**
	 * @name Scanner stuff
	 * @{
	 */

	/** Have we started to scan the input stream? */
	stream_start_produced bool

	/** Have we reached the end of the input stream? */
	stream_end_produced bool

	/** The number of unclosed '[' and '{' indicators. */
	flow_level int

	/** The tokens queue. */
	tokens      []yaml_token_t
	tokens_head int

	/** The number of tokens fetched from the queue. */
	tokens_parsed int

	/* Does the tokens queue contain a token ready for dequeueing. */
	token_available bool

	/** The indentation levels stack. */
	indents []int

	/** The current indentation level. */
	indent int

	/** May a simple key occur at the current position? */
	simple_key_allowed bool

	/** The stack of simple keys. */
	simple_keys []yaml_simple_key_t

	/**
	 * @}
	 */

	/**
	 * @name Parser stuff
	 * @{
	 */

	/** The parser states stack. */
	states []yaml_parser_state_t

	/** The current parser state. */
	state yaml_parser_state_t

	/** The stack of marks. */
	marks []yaml_mark_t

	/** The list of TAG directives. */
	tag_directives []yaml_tag_directive_t

	/**
	 * @}
	 */

	/**
	 * @name Dumper stuff
	 * @{
	 */

	/** The alias data. */
	aliases []yaml_alias_data_t

	/** The currently parsed document. */
	document *yaml_document_t

	/**
	 * @}
	 */

}

/**
 * The prototype of a write handler.
 *
 * The write handler is called when the emitter needs to flush the accumulated
 * characters to the output.  The handler should write @a size bytes of the
 * @a buffer to the output.
 *
 * @param[in,out]   data        A pointer to an application data specified by
 *                              yaml_emitter_set_output().
 * @param[in]       buffer      The buffer with bytes to be written.
 * @param[in]       size        The size of the buffer.
 *
 * @returns On success, the handler should return @c 1.  If the handler failed,
 * the returned value should be @c 0.
 */

type yaml_write_handler_t func(emitter *yaml_emitter_t, buffer []byte) error

/** The emitter states. */
type yaml_emitter_state_t int

const (
	/** Expect STREAM-START. */
	YAML_EMIT_STREAM_START_STATE yaml_emitter_state_t = iota
	/** Expect the first DOCUMENT-START or STREAM-END. */
	YAML_EMIT_FIRST_DOCUMENT_START_STATE
	/** Expect DOCUMENT-START or STREAM-END. */
	YAML_EMIT_DOCUMENT_START_STATE
	/** Expect the content of a document. */
	YAML_EMIT_DOCUMENT_CONTENT_STATE
	/** Expect DOCUMENT-END. */
	YAML_EMIT_DOCUMENT_END_STATE
	/** Expect the first item of a flow sequence. */
	YAML_EMIT_FLOW_SEQUENCE_FIRST_ITEM_STATE
	/** Expect an item of a flow sequence. */
	YAML_EMIT_FLOW_SEQUENCE_ITEM_STATE
	/** Expect the first key of a flow mapping. */
	YAML_EMIT_FLOW_MAPPING_FIRST_KEY_STATE
	/** Expect a key of a flow mapping. */
	YAML_EMIT_FLOW_MAPPING_KEY_STATE
	/** Expect a value for a simple key of a flow mapping. */
	YAML_EMIT_FLOW_MAPPING_SIMPLE_VALUE_STATE
	/** Expect a value of a flow mapping. */
	YAML_EMIT_FLOW_MAPPING_VALUE_STATE
	/** Expect the first item of a block sequence. */
	YAML_EMIT_BLOCK_SEQUENCE_FIRST_ITEM_STATE
	/** Expect an item of a block sequence. */
	YAML_EMIT_BLOCK_SEQUENCE_ITEM_STATE
	/** Expect the first key of a block mapping. */
	YAML_EMIT_BLOCK_MAPPING_FIRST_KEY_STATE
	/** Expect the key of a block mapping. */
	YAML_EMIT_BLOCK_MAPPING_KEY_STATE
	/** Expect a value for a simple key of a block mapping. */
	YAML_EMIT_BLOCK_MAPPING_SIMPLE_VALUE_STATE
	/** Expect a value of a block mapping. */
	YAML_EMIT_BLOCK_MAPPING_VALUE_STATE
	/** Expect nothing. */
	YAML_EMIT_END_STATE
)

/**
 * The emitter structure.
 *
 * All members are internal.  Manage the structure using the @c yaml_emitter_
 * family of functions.
 */

type yaml_emitter_t struct {

	/**
	 * @name Error handling
	 * @{
	 */

	/** Error type. */
	error yaml_error_type_t
	/** Error description. */
	problem string

	/**
	 * @}
	 */

	/**
	 * @name Writer stuff
	 * @{
	 */

	/** Write handler. */
	write_handler yaml_write_handler_t

	/** Standard (string or file) output data. */
	output_buffer *[]byte
	output_file   io.Writer

	/** The working buffer. */
	buffer     []byte
	buffer_pos int

	/** The raw buffer. */
	raw_buffer     []byte
	raw_buffer_pos int

	/** The stream encoding. */
	encoding yaml_encoding_t

	/**
	 * @}
	 */

	/**
	 * @name Emitter stuff
	 * @{
	 */

	/** If the output is in the canonical style? */
	canonical bool
	/** The number of indentation spaces. */
	best_indent int
	/** The preferred width of the output lines. */
	best_width int
	/** Allow unescaped non-ASCII characters? */
	unicode bool
	/** The preferred line break. */
	line_break yaml_break_t

	/** The stack of states. */
	states []yaml_emitter_state_t

	/** The current emitter state. */
	state yaml_emitter_state_t

	/** The event queue. */
	events      []yaml_event_t
	events_head int

	/** The stack of indentation levels. */
	indents []int

	/** The list of tag directives. */
	tag_directives []yaml_tag_directive_t

	/** The current indentation level. */
	indent int

	/** The current flow level. */
	flow_level int

	/** Is it the document root context? */
	root_context bool
	/** Is it a sequence context? */
	sequence_context bool
	/** Is it a mapping context? */
	mapping_context bool
	/** Is it a simple mapping key context? */
	simple_key_context bool

	/** The current line. */
	line int
	/** The current column. */
	column int
	/** If the last character was a whitespace? */
	whitespace bool
	/** If the last character was an indentation character (' ', '-', '?', ':')? */
	indention bool
	/** If an explicit document end is required? */
	open_ended bool

	/** Anchor analysis. */
	anchor_data struct {
		/** The anchor value. */
		anchor []byte
		/** Is it an alias? */
		alias bool
	}

	/** Tag analysis. */
	tag_data struct {
		/** The tag handle. */
		handle []byte
		/** The tag suffix. */
		suffix []byte
	}

	/** Scalar analysis. */
	scalar_data struct {
		/** The scalar value. */
		value []byte
		/** Does the scalar contain line breaks? */
		multiline bool
		/** Can the scalar be expessed in the flow plain style? */
		flow_plain_allowed bool
		/** Can the scalar be expressed in the block plain style? */
		block_plain_allowed bool
		/** Can the scalar be expressed in the single quoted style? */
		single_quoted_allowed bool
		/** Can the scalar be expressed in the literal or folded styles? */
		block_allowed bool
		/** The output style. */
		style yaml_scalar_style_t
	}

	/**
	 * @}
	 */

	/**
	 * @name Dumper stuff
	 * @{
	 */

	/** If the stream was already opened? */
	opened bool
	/** If the stream was already closed? */
	closed bool

	/** The information associated with the document nodes. */
	anchors *struct {
		/** The number of references. */
		references int
		/** The anchor id. */
		anchor int
		/** If the node has been emitted? */
		serialized bool
	}

	/** The last assigned anchor id. */
	last_anchor_id int

	/** The currently emitted document. */
	document *yaml_document_t

	/**
	 * @}
	 */

}
