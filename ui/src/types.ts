// DSS Type Definitions for the UI

export interface DesignSystem {
  meta: Meta;
  principles?: Principle[];
  foundations: Foundations;
  components?: Component[];
  patterns?: Pattern[];
  templates?: Template[];
  content?: Content;
  accessibility?: Accessibility;
  governance?: Governance;
  themeBindings?: ThemeBindings[];
  modes?: string[];
}

export interface Meta {
  name: string;
  version: string;
  description?: string;
  organization?: string;
  repository?: string;
  documentationUrl?: string;
  figmaUrl?: string;
  lastUpdated?: string;
  maintainers?: string[];
}

export interface Principle {
  id: string;
  name: string;
  description?: string;
}

export interface Foundations {
  colors?: ColorToken[];
  typography?: Typography;
  spacing?: SpacingScale;
  elevation?: ElevationToken[];
  motion?: MotionSystem;
  grid?: GridSystem;
  breakpoints?: Breakpoint[];
  borderRadius?: BorderRadiusToken[];
  borderWidth?: BorderWidthToken[];
  opacity?: OpacityToken[];
  zIndex?: ZIndexToken[];
  densities?: DensityToken[];
}

export interface DensityToken {
  id: string;
  scale: number;
  description?: string;
  spacingOverrides?: Record<string, string>;
}

export interface ColorToken {
  id: string;
  value: string;
  semantic?: string;
  usage?: string;
  contrast?: Contrast;
  lightModeValue?: string;
  darkModeValue?: string;
  modes?: Record<string, string>;
}

export interface Contrast {
  onLight?: string;
  onDark?: string;
  ratio?: number;
  wcagLevel?: 'AA' | 'AAA';
}

export interface Typography {
  fontFamilies?: FontFamily[];
  fontSizes?: FontSize[];
  fontWeights?: FontWeight[];
  lineHeights?: LineHeight[];
  letterSpacings?: LetterSpacing[];
  typeScale?: TypeStyle[];
}

export interface FontFamily {
  id: string;
  name: string;
  stack?: string;
  value?: string;
  weights?: number[];
  source?: string;
  usage?: string;
}

export interface FontSize {
  id: string;
  value: string;
  pixelValue?: number;
}

export interface FontWeight {
  id: string;
  value: number;
}

export interface LineHeight {
  id: string;
  value: string;
}

export interface LetterSpacing {
  id: string;
  value: string;
}

export interface TypeStyle {
  id: string;
  name: string;
  fontFamily: string;
  fontSize: string;
  fontWeight: string;
  lineHeight: string;
  letterSpacing?: string;
  usage?: string;
}

export interface SpacingScale {
  baseUnit: string;
  scale: SpacingToken[];
}

export interface SpacingToken {
  id: string;
  value: string;
  pixelValue?: number;
}

export interface ElevationToken {
  id: string;
  value: string;
  usage?: string;
}

export interface MotionSystem {
  durations?: DurationToken[];
  easings?: EasingToken[];
}

export interface DurationToken {
  id: string;
  value: string;
  usage?: string;
}

export interface EasingToken {
  id: string;
  value: string;
  usage?: string;
}

export interface GridSystem {
  columns: number;
  gutterWidth: string;
  marginWidth: string;
  maxWidth?: string;
}

export interface Breakpoint {
  id: string;
  minWidth: string;
  maxWidth?: string;
  columns?: number;
}

export interface BorderRadiusToken {
  id: string;
  value: string;
}

export interface BorderWidthToken {
  id: string;
  value: string;
}

export interface OpacityToken {
  id: string;
  value: string;
}

export interface ZIndexToken {
  id: string;
  value: number;
  usage?: string;
}

export interface Component {
  id: string;
  name: string;
  description?: string;
  category?: string;
  variants?: Variant[];
  states?: State[];
  props?: Prop[];
  events?: ComponentEvent[];
  slots?: Slot[];
  uses?: string[];
  tokensUsed?: string[];
  constraints?: Constraints;
  accessibility?: ComponentA11y;
  llm?: LLMContext;
  themingContract?: ThemingContract;
  documentationUrl?: string;
  figmaUrl?: string;
}

export interface Variant {
  id: string;
  name: string;
  description?: string;
  isDefault?: boolean;
  tokenOverrides?: TokenOverride[];
}

export interface TokenOverride {
  tokenId: string;
  value: string;
}

export interface State {
  id: string;
  name?: string;
  description?: string;
  tokenOverrides?: TokenOverride[];
}

export interface Prop {
  name: string;
  type: string;
  description?: string;
  required?: boolean;
  default?: unknown;
  enumValues?: string[];
  constraints?: PropConstraints;
}

export interface PropConstraints {
  min?: number;
  max?: number;
  minLength?: number;
  maxLength?: number;
  pattern?: string;
  step?: number;
}

export interface ComponentEvent {
  name: string;
  description?: string;
  bubbles?: boolean;
  composed?: boolean;
  cancelable?: boolean;
  detail?: EventDetailField[];
}

export interface EventDetailField {
  name: string;
  type: string;
  description?: string;
}

export interface Slot {
  name: string;
  description?: string;
  required?: boolean;
  allowedComponents?: string[];
}

export interface Constraints {
  minWidth?: string;
  maxWidth?: string;
  maxChildren?: number;
  requiredParent?: string;
  forbiddenChildren?: string[];
  responsiveBehavior?: string;
}

export interface ComponentA11y {
  role?: string;
  requiredAttributes?: string[];
  keyboardSupport?: KeyboardInteraction[];
  screenReaderNotes?: string;
  focusManagement?: string;
}

export interface KeyboardInteraction {
  key: string;
  action: string;
}

export interface LLMContext {
  intent?: string;
  allowedContexts?: string[];
  forbiddenContexts?: string[];
  exampleUsage?: string[];
  antiPatterns?: string[];
  semanticMeaning?: string;
  relatedElements?: string[];
  priorityScore?: number;
}

export interface ThemingContract {
  prefix: string;
  description?: string;
  tokens: ThemeToken[];
}

export interface ThemeToken {
  id: string;
  cssProperty: string;
  semantic?: string;
  description?: string;
  defaultLight?: string;
  defaultDark?: string;
  defaults?: Record<string, string>;
  densitySensitive?: boolean;
}

export interface ThemeBindings {
  component: string;
  specUrl?: string;
  themeMode?: string;
  strategy?: 'explicit' | 'semantic' | 'inherit';
  mappings: TokenMapping[];
}

export interface TokenMapping {
  from: string;
  to: string;
  transform?: string;
}

export interface Pattern {
  id: string;
  name: string;
  description?: string;
}

export interface Template {
  id: string;
  name: string;
  description?: string;
}

export interface Content {
  voiceTone?: string;
  guidelines?: string[];
}

export interface Accessibility {
  wcagLevel?: string;
  requirements?: string[];
}

export interface Governance {
  versioning?: string;
  deprecationPolicy?: string;
}
