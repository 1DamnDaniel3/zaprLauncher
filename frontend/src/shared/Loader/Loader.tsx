import s from './Loader.module.css'

export function Loader() {
    return (
        <div className={s.wrapper}>
            <div className={s.loader} />
        </div>
    )
}