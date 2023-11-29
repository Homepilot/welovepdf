import './Button.css';

type ButtonProps = {
    clicked(): void | Promise<void>
    disabled?: boolean;
    extraClasses?: string[]
}

export const Button: React.FC<React.PropsWithChildren & ButtonProps> = ({
    clicked,
    children,
    disabled = false,
    extraClasses = []
}) => (
    <span 
        onClick={() => {
            if(disabled) return
            clicked()
        }}
        className={[
            'hp-btn',
            disabled ? 'hp-btn-disabled' : 'hp-btn-enabled',
            ...extraClasses
        ].filter(Boolean).join(' ')}
    >
        {children}
    </span>
                        
)