import { API_URL } from "@/utils/config";
import axios from "axios";

class APIClient {
  /** @type {string} */
  url = API_URL;

  /**
   * @param {string} cik
   * @param {string} accNum Accession number.
   * @returns {Promise<Array<Object>>} Balance sheet objects.
   */
  async getBalanceSheet(cik, accNum) {
    const { data, status } = await axios.get(
      `${this.url}/sheets/${cik}/${accNum}`
    );
    if (status != 200) return statusError(status);
    return data;
  }

  /**
   * @param {string} cik
   * @param {string} accNum Accession number.
   * @returns {Promise<Array<Object>>} Finance note objects.
   */
  async getFinancialNotes(cik, accNum) {
    const { data, status } = await axios.get(
      `${this.url}/notes/${cik}/${accNum}`
    );
    if (status != 200) return statusError(status);
    return data;
  }
}

function statusError(status) {
  return { error: `Server returned invalid status (${status})` };
}

export default APIClient;
