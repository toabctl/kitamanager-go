'use client';

import * as React from 'react';
import { X, Plus } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';

export interface TagInputProps {
  value: string[];
  onChange: (value: string[]) => void;
  placeholder?: string;
  disabled?: boolean;
  className?: string;
  id?: string;
  suggestions?: string[];
  suggestionsLabel?: string;
  /** Map of attribute name to its exclusive group. Tags in the same group are mutually exclusive. */
  exclusiveGroupMap?: Record<string, string | null>;
}

export function TagInput({
  value,
  onChange,
  placeholder = 'Type and press Enter...',
  disabled = false,
  className,
  id,
  suggestions = [],
  suggestionsLabel,
  exclusiveGroupMap = {},
}: TagInputProps) {
  const [inputValue, setInputValue] = React.useState('');
  const inputRef = React.useRef<HTMLInputElement>(null);

  const addTag = (tag: string) => {
    const trimmed = tag.trim().toLowerCase();
    if (!trimmed || value.includes(trimmed)) {
      setInputValue('');
      return;
    }

    // Check if the new tag has an exclusive group
    const newTagGroup = exclusiveGroupMap[trimmed];

    let newValue = [...value];

    // If it has an exclusive group, remove other tags from the same group
    if (newTagGroup) {
      newValue = newValue.filter((existingTag) => {
        const existingGroup = exclusiveGroupMap[existingTag];
        return existingGroup !== newTagGroup;
      });
    }

    newValue.push(trimmed);
    onChange(newValue);
    setInputValue('');
  };

  const removeTag = (tagToRemove: string) => {
    onChange(value.filter((tag) => tag !== tagToRemove));
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' || e.key === ',') {
      e.preventDefault();
      if (inputValue.trim()) {
        addTag(inputValue);
      }
    } else if (e.key === 'Backspace' && !inputValue && value.length > 0) {
      removeTag(value[value.length - 1]);
    }
  };

  const handleBlur = () => {
    if (inputValue.trim()) {
      addTag(inputValue);
    }
  };

  // Filter suggestions to only show ones not already selected
  const availableSuggestions = suggestions.filter((s) => !value.includes(s.toLowerCase()));

  // Get the exclusive group for a tag (for visual indication)
  const getTagGroup = (tag: string): string | null => {
    return exclusiveGroupMap[tag] || null;
  };

  return (
    <div className="space-y-2">
      <div
        className={cn(
          'flex min-h-10 w-full flex-wrap items-center gap-1.5 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-within:ring-2 focus-within:ring-ring focus-within:ring-offset-2',
          disabled && 'cursor-not-allowed opacity-50',
          className
        )}
        onClick={() => inputRef.current?.focus()}
      >
        {value.map((tag) => (
          <Badge
            key={tag}
            variant={getTagGroup(tag) ? 'default' : 'secondary'}
            className="gap-1 pr-1"
            title={getTagGroup(tag) ? `Group: ${getTagGroup(tag)}` : undefined}
          >
            {tag}
            {!disabled && (
              <button
                type="button"
                onClick={(e) => {
                  e.stopPropagation();
                  removeTag(tag);
                }}
                className="ml-1 rounded-full outline-none ring-offset-background hover:bg-muted focus:ring-2 focus:ring-ring focus:ring-offset-2"
                aria-label={`Remove ${tag}`}
              >
                <X className="h-3 w-3" />
              </button>
            )}
          </Badge>
        ))}
        <input
          ref={inputRef}
          id={id}
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyDown={handleKeyDown}
          onBlur={handleBlur}
          placeholder={value.length === 0 ? placeholder : ''}
          disabled={disabled}
          className="flex-1 bg-transparent outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed"
          style={{ minWidth: '80px' }}
        />
      </div>
      {availableSuggestions.length > 0 && !disabled && (
        <div className="flex flex-wrap gap-1">
          {suggestionsLabel && (
            <span className="mr-1 self-center text-xs text-muted-foreground">
              {suggestionsLabel}
            </span>
          )}
          {availableSuggestions.map((suggestion) => {
            const group = exclusiveGroupMap[suggestion];
            return (
              <button
                key={suggestion}
                type="button"
                onClick={() => addTag(suggestion)}
                className={cn(
                  'inline-flex items-center gap-1 rounded-full border border-dashed px-2 py-0.5 text-xs transition-colors',
                  group
                    ? 'border-primary/50 text-primary/70 hover:border-primary hover:text-primary'
                    : 'border-muted-foreground/50 text-muted-foreground hover:border-primary hover:text-primary'
                )}
                title={group ? `Group: ${group} (selecting replaces others in group)` : undefined}
              >
                <Plus className="h-3 w-3" />
                {suggestion}
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}
