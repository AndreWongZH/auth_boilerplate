export const Field = ({ title, placeholder = '', type = 'text', value, onChange }:
    { title: string, placeholder?: string, type?: string, value: string, onChange: Function }) => {
    return (
        <div>
            <h3>{title}</h3>
            <input className="py-2 px-4 leading-tight text-gray-700 appearance-none border-4 border-gray-200 rounded focus:outline-none focus:bg-white focus:border-teal-500"
                placeholder={placeholder} type={type} value={value} onChange={(e) => onChange(e.target.value)} />
        </div>
    )
}

export const Button = ({ text, onClick }: { text: string, onClick?: Function }) => {
    return (
       <button onClick={() => { onClick?.() }} 
            className="px-7 py-3 mt-10 bg-teal-500 hover:bg-teal-700 border rounded-lg font-bold">{text}</button>
    )
}

export const ButtonText = ({ text, onClick }: { text: string, onClick: Function }) => {
    return (
        <button onClick={() => { onClick() }} 
            className="px-7 py-3 mt-10 hover:bg-gray-600 rounded-lg font-bold">{text}</button>
    )
}
