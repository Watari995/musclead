import 'package:flutter/material.dart';

import '../theme/app_tokens.dart';

/// −/＋ で増減できる数値入力（キーボード不要）。直接タイプも可能。
class NumberStepper extends StatelessWidget {
  const NumberStepper({
    super.key,
    required this.label,
    required this.controller,
    this.step = 0.1,
    this.min = 0,
    this.max,
    this.hint,
  });

  final String label;
  final TextEditingController controller;
  final double step;
  final double min;
  final double? max;
  final String? hint;

  int get _decimals {
    final s = step.toString();
    final i = s.indexOf('.');
    return i < 0 ? 0 : s.length - i - 1;
  }

  void _nudge(double delta) {
    final cur = double.tryParse(controller.text.trim()) ?? min;
    var next = cur + delta;
    if (next < min) next = min;
    if (max != null && next > max!) next = max!;
    controller.text = next.toStringAsFixed(_decimals);
  }

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.only(left: 2, bottom: 6),
          child: Text(
            label,
            style: const TextStyle(fontSize: 13, fontWeight: FontWeight.w600),
          ),
        ),
        Container(
          decoration: BoxDecoration(
            border: Border.all(color: t.border),
            borderRadius: BorderRadius.circular(13),
          ),
          child: Row(
            children: [
              _StepButton(icon: Icons.remove, onTap: () => _nudge(-step)),
              Expanded(
                child: TextField(
                  controller: controller,
                  textAlign: TextAlign.center,
                  keyboardType: const TextInputType.numberWithOptions(
                    decimal: true,
                  ),
                  style: const TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w700,
                  ),
                  decoration: InputDecoration(
                    hintText: hint,
                    hintStyle: TextStyle(
                      color: t.subtle,
                      fontWeight: FontWeight.w400,
                    ),
                    border: InputBorder.none,
                    isDense: true,
                  ),
                ),
              ),
              _StepButton(icon: Icons.add, onTap: () => _nudge(step)),
            ],
          ),
        ),
      ],
    );
  }
}

class _StepButton extends StatelessWidget {
  const _StepButton({required this.icon, required this.onTap});

  final IconData icon;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        width: 52,
        height: 52,
        alignment: Alignment.center,
        child: Icon(icon, color: t.accent, size: 22),
      ),
    );
  }
}
