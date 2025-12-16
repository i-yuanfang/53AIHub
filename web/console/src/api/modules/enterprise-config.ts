import service from '../config'
import { handleError } from '../errorHandler'

export const enterpriseConfigApi = {
    getTypes() {
        return service.get('/api/enterprise-configs').catch(handleError)
    },
    get(type: string) {
        return service.get(`/api/enterprise-configs/${type}`).catch(handleError)
    },
    save(type: string, data: { content: string; enabled: boolean }) {
        return service.post(`/api/enterprise-configs/${type}`, data).catch(handleError)
    },
    toggle(type: string) {
        return service.put(`/api/enterprise-configs/${type}/toggle`).catch(handleError)
    },
    isEnabled(type: string) {
        return service.get(`/api/enterprise-configs/${type}/enabled`).catch(handleError)
    }
}

export default enterpriseConfigApi
