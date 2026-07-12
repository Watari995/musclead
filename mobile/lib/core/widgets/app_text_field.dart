import 'package:flutter/material.dart';

import '../theme/app_tokens.dart';
import '../theme/sketchy.dart';

/// プレビュー準拠の入力欄。ラベル + 手描き輪郭 + 16px(iOS の自動ズーム回避)。
class AppTextField extends StatefulWidget {
  const AppTextField({
    super.key,
    required this.label,
    this.controller,
    this.hint,
    this.obscureText = false,
    this.keyboardType,
    this.textInputAction,
    this.errorText,
    this.onChanged,
    this.onSubmitted,
    this.focusNode,
    this.enabled = true,
    this.autofillHints,
  });

  final String label;
  final TextEditingController? controller;
  final String? hint;
  final bool obscureText;
  final TextInputType? keyboardType;
  final TextInputAction? textInputAction;
  final String? errorText;
  final ValueChanged<String>? onChanged;
  final ValueChanged<String>? onSubmitted;
  final FocusNode? focusNode;
  final bool enabled;
  final Iterable<String>? autofillHints;

  @override
  State<AppTextField> createState() => _AppTextFieldState();
}

class _AppTextFieldState extends State<AppTextField> {
  late final FocusNode _focusNode = widget.focusNode ?? FocusNode();
  late final bool _ownsFocusNode = widget.focusNode == null;
  bool _focused = false;

  @override
  void initState() {
    super.initState();
    _focusNode.addListener(_handleFocusChange);
  }

  void _handleFocusChange() {
    setState(() => _focused = _focusNode.hasFocus);
  }

  @override
  void dispose() {
    _focusNode.removeListener(_handleFocusChange);
    if (_ownsFocusNode) _focusNode.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final hasError = widget.errorText != null;
    final strokeColor = hasError
        ? context.colors.error
        : _focused
        ? t.accent
        : t.border;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.only(left: 2, bottom: 6),
          child: Text(
            widget.label,
            style: TextStyle(fontFamily: 'Caveat', fontSize: 17, color: t.ink),
          ),
        ),
        RoughBox(
          color: strokeColor,
          radius: BorderRadius.circular(13),
          padding: const EdgeInsets.symmetric(horizontal: 13),
          child: TextField(
            controller: widget.controller,
            focusNode: _focusNode,
            obscureText: widget.obscureText,
            keyboardType: widget.keyboardType,
            textInputAction: widget.textInputAction,
            onChanged: widget.onChanged,
            onSubmitted: widget.onSubmitted,
            enabled: widget.enabled,
            autofillHints: widget.autofillHints,
            style: const TextStyle(fontSize: 16),
            decoration: InputDecoration(
              hintText: widget.hint,
              hintStyle: TextStyle(color: t.subtle, fontSize: 16),
              border: InputBorder.none,
              isDense: true,
              contentPadding: const EdgeInsets.symmetric(vertical: 14),
            ),
          ),
        ),
        if (hasError)
          Padding(
            padding: const EdgeInsets.only(left: 2, top: 6),
            child: Text(
              widget.errorText!,
              style: TextStyle(color: context.colors.error, fontSize: 12.5),
            ),
          ),
      ],
    );
  }
}
